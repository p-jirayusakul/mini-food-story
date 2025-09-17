package redis

import (
	"context"
	"errors"
	"food-story/pkg/exceptions"
	"log"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisInterface interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
	TTL(key string) (time.Duration, error)
	Close()
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

func (r *RedisClient) Set(key string, value string, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	data, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		slog.Error("Redis key not found: ", "key", key)
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

func (r *RedisClient) TTL(key string) (time.Duration, error) {
	data, err := r.Client.TTL(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		slog.Error("Redis key not found: ", "key", key)
		return 0, exceptions.ErrRedisKeyNotFound
	} else if err != nil {
		return 0, err
	} else {
		return data, nil
	}
}

func (r *RedisClient) Close() {
	if r.Client == nil {
		log.Printf("RedisClient is nil. Skipping Close.")
		return
	}
	err := r.Client.Close()
	if err != nil {
		log.Printf("Unable to close Redis: %v", err)
	}
}
