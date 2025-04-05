package consul

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/saas-flow/monorepo/libs/config"
	"go.uber.org/fx"
)

var Register = fx.Module("consule.register", fx.Invoke(RegisterService))

func RegisterService() {
	consulConfig := api.DefaultConfig()

	serviceProtocol := strings.ToLower(config.GetString("SERVICE_PROTOCOL"))

	if config.GetBool("CONSUL.TLS.ENABLED") {
		tlsConfig := &tls.Config{RootCAs: nil} // Bisa ditambahkan CA jika diperlukan
		consulConfig.TLSConfig = api.TLSConfig{CAFile: config.GetString("CONSUL.TLS.CA_CERT")}
		consulConfig.HttpClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Failed to connect to Consul: %v", err)
	}

	serviceName := config.GetString("SERVICE_NAME")
	servicePort := config.GetInt("SERVICE_PORT")
	serviceID := fmt.Sprintf("%s-%d", serviceName, servicePort)

	// Definisi Service
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceName,
		Port:    servicePort,
		Meta: map[string]string{
			"protocol": serviceProtocol,
		},
		Check: &api.AgentServiceCheck{
			HTTP:          fmt.Sprintf("%s://%s:%d/health", serviceProtocol, serviceName, servicePort),
			Interval:      config.GetString("CONSUL.SERVICE_CHECK.INTERVAL"),
			TLSSkipVerify: serviceProtocol == "https",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	log.Printf("Service registered: %s", serviceID)
}

func DiscoverService(serviceName string) (string, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return "", fmt.Errorf("failed to connect to Consul: %w", err)
	}

	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil || len(services) == 0 {
		return "", fmt.Errorf("service %s not found", serviceName)
	}

	// Ambil instance pertama yang sehat
	service := services[0].Service
	protocol := service.Meta["protocol"]
	return fmt.Sprintf("%s://%s:%d", protocol, service.Address, service.Port), nil
}
