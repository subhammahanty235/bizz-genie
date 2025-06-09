package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type InternalRedisClient struct {
	client *redis.Client
}

func NewInternalRedisClient() (*InternalRedisClient, error) {

	// var ops *redis.Options
	opts, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	internalClient := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, errcon := internalClient.Ping(ctx).Result()
	if errcon != nil {
		return nil, fmt.Errorf("failed to connect to Internal Redis Instance: %w", err)
	}

	fmt.Println(status)

	fmt.Println("🔌 Connected to Internal BizzMQ client database")

	return &InternalRedisClient{
		client: internalClient,
	}, nil

}

func (r *InternalRedisClient) GetInternalRedisClient() *redis.Client {
	return r.client
}

func (r *InternalRedisClient) CloseInternalRedisClient() error {
	return r.client.Close()
}
