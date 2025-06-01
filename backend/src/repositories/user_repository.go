package repositories

import (
	"backend/src/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUserRepository struct {
	client *redis.Client
	ctx		context.Context
}

func NewRedisUserRepository(client *redis.Client) *RedisUserRepository {
	return &RedisUserRepository{
		client: client,
		ctx: context.Background(),
	}
}

func (r *RedisUserRepository) CreateUser(user model.User) (uint64, error) {
	if user.ID == 0 {
		newID, err := r.client.Incr(r.ctx, "user:counter").Result()
		if err != nil {
			return 0, fmt.Errorf("failed to generate user ID: %w", err)
		}
		user.ID = uint64(newID)
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal user data: %w", err)
	}

	key := fmt.Sprintf("user:%d", user.ID)

	err = r.client.Set(r.ctx, key, userData, 24*time.Hour).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to store user in redis: %w", err)
	}

	return user.ID, nil
}