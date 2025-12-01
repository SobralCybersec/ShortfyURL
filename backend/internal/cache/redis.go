package cache

import (
	"context"
	"encoding/json"
	"shortfyurl/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisCache) Set(key string, value *models.URL, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, expiration).Err()
}

func (r *RedisCache) Get(key string) (*models.URL, error) {
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var url models.URL
	if err := json.Unmarshal([]byte(data), &url); err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
