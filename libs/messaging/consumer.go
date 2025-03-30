package messaging

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/saas-flow/monorepo/libs/config"
	"go.uber.org/fx"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer() (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetString("KAFKA.BROKERS"),
		"group.id":          config.GetString("KAFKA.GROUP_ID"),
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: c}, nil
}

func (c *Consumer) Subscribe(topics []string) error {
	return c.consumer.SubscribeTopics(topics, nil)
}

func (c *Consumer) Consume(ctx context.Context, handler func(msg *kafka.Message)) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping consumer...")
			c.consumer.Close()
			return
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err == nil {
				handler(msg)
			} else {
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}
}

var ConsumerModule = fx.Module("consumer",
	fx.Provide(
		NewConsumer,
	),
)
