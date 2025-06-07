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
	ops := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	}

	internalClient := redis.NewClient(ops)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := internalClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Internal Redis Instance: %w", err)
	}

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
