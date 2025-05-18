package redisclient

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type externalRedisClient struct {
	client *redis.Client
}

func NewExternalRedisClient() (*externalRedisClient, error) {
	redisHost := os.Getenv("REDIS-HOST")
	redisPort := os.Getenv("REDIS-PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisUrl := os.Getenv("REDIS_URL")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	// There will be 2 scenarios, one will be using the host, port and password and another with url
	var ops *redis.Options
	if redisUrl != "" {
		ops = &redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
		}
	} else {
		opts, err := redis.ParseURL(redisUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
		}
		ops = opts
	}

	externalClient := redis.NewClient(ops)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := externalClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("🔌 Connected to External BizzMQ client database")

	return &externalRedisClient{
		client: externalClient,
	}, nil

}

func (r *externalRedisClient) GetExternalRedisClient() *redis.Client {
	return r.client
}

func (r *externalRedisClient) CloseRedisClient() error {
	return r.client.Close()
}
