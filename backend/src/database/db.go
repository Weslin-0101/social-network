package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		if parsed, err := strconv.Atoi(dbStr); err == nil {
			db = parsed
		}
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
	})

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Printf("Connected to Redis: %s\n", pong)

	return client, nil
}