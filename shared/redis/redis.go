package redis

import (
	"context"
	"errors"
	"food-story/pkg/exceptions"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
)

var ctx = context.Background()

type RedisInterface interface {
	Set(key string, value string, expiration int) error
	Get(key string) (string, error)
	Del(key string) error
}

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string, password string, db int) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,     // Redis server address
		Password: password, // Password (empty if no password)
		DB:       db,       // Default database
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
	}

	return &RedisClient{Client: client}
}

var _ RedisInterface = (*RedisClient)(nil)

func (r *RedisClient) Set(key string, value string, expiration int) error {
	return r.Client.Set(ctx, key, value, 0).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	data, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		slog.Error("Key not found", "key", key)
		return "", exceptions.ErrRedisKeyNotFound
	} else if err != nil {
		return "", err
	} else {
		return data, nil
	}
}

func (r *RedisClient) Del(key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisClient) close() {
	r.Client.Close()
}
