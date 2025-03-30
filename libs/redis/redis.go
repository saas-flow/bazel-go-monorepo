package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var Module = fx.Module("redis",
	fx.Provide(NewRedisClient),
)

func NewRedisClient(opts *redis.Options) (*redis.Client, error) {

	client := redis.NewClient(opts)

	// Tes koneksi
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("failed connect to redis: %v", err)
		return nil, err
	}

	fmt.Println("Redis connected!")
	return client, nil
}
