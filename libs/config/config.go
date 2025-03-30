package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

func init() {

	// configPath := os.Getenv("CONFIG_PATH")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	provider := os.Getenv("REMOTE_CONFIG_PROVIDER")
	if provider != "" {
		endpoint := os.Getenv("REMOTE_CONFIG_PROVIDER_ENDPOINT")
		path := os.Getenv("REMOTE_CONFIG_PROVIDER_PATH")
		if err := viper.AddRemoteProvider(
			provider,
			endpoint,
			path,
		); err != nil {
			log.Fatalf("failed viper.AddRemoteProvider: %v\n", err)
			panic(err)
		}

		if err := viper.ReadRemoteConfig(); err != nil {
			log.Fatalf("failed viper.ReadRemoteConfig: %v\n", err)
			panic(err)
		}
	} else {
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("failed viper.ReadInConfig: %v\n", err)
			panic(err)
		}
	}
}

// GetString
//
// Parameters:
//   - key: string
//
// Result
//   - string:
func GetString(key string) string {
	return viper.GetString(key)
}

// GetStringSlice
//
// Parameters:
//   - key: string
//
// Result
//   - []string:
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// GetStringMap
//
// Parameters:
//   - key: string
//
// Result
//   - map[string]any:
func GetStringMap(key string) map[string]any {
	return viper.GetStringMap(key)
}

// GetStringMapString
//
// Parameters:
//   - key: string
//
// Result
//   - map[string]string:
func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

// GetStringMapStringSlice
//
// Parameters:
//   - key: string
//
// Result
//   - map[string][]string:
func GetStringMapStringSlice(key string) map[string][]string {
	return viper.GetStringMapStringSlice(key)
}

// GetBool
//
// Parameters:
//   - key: string
//
// Result
//   - bool:
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetInt
//
// Parameters:
//   - key: string
//
// Result
//   - int:
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetIntSlice
//
// Parameters:
//   - key: string
//
// Result
//   - []int:
func GetIntSlice(key string) []int {
	return viper.GetIntSlice(key)
}

// GetInt32
//
// Parameters:
//   - key: string
//
// Result
//   - int32:
func GetInt32(key string) int32 {
	return viper.GetInt32(key)
}

// GetInt64
//
// Parameters:
//   - key: string
//
// Result
//   - int64:
func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

// GetFloat64
//
// Parameters:
//   - key: string
//
// Result
//   - float64:
func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

// GetDuration
//
// Parameters:
//   - key: string
//
// Result
//   - time.Duration:
func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
