package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type internalRedisClient struct {
	client *redis.Client
}

func NewInternalRedisClient() (*internalRedisClient, error) {

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

	return &internalRedisClient{
		client: internalClient,
	}, nil

}

func (r *internalRedisClient) GetInternalRedisClient() *redis.Client {
	return r.client
}

func (r *internalRedisClient) CloseInternalRedisClient() error {
	return r.client.Close()
}
