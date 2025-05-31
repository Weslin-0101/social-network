package redis

import (
	"backend/src/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func PopulateUsers() error {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
		Password: "",
		DB: 0,
	})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis successfully")

	users := []model.User{
		{
			ID:	   101,
			Username: "john_doe",
			Nickname: "John",
			Email:  "john@gmail.com",
			Password: "hash_password",
			CreatedAt: time.Now(),
		},
	}

	for _, user := range users {
		userData, err := json.Marshal(user)
		if err != nil {
			log.Printf("could not marshal user %v: %v", user, err)
			continue
		}

		key := fmt.Sprintf("user:%d", user.ID)

		err = rdb.Set(ctx, key, userData, 0).Err()
		if err != nil {
			log.Printf("could not set user %d in Redis: %v", user.ID, err)
			continue
		}
	}

	fmt.Println("Users populated successfully")
	return nil
}