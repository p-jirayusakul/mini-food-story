package cache

import (
	"encoding/json"
	"food-story/shared/redis"
	"food-story/table-service/internal/domain"
	"time"
)

type RedisTableCacheInterface interface {
	GetCachedTable(key string) (*domain.CurrentTableSession, error)
	SetCachedTable(key string, table *domain.CurrentTableSession, ttl time.Duration) error
	DeleteCachedTable(key string) error
}

type RedisTableCache struct {
	client redis.RedisInterface
}

func NewRedisTableCache(client *redis.RedisClient) *RedisTableCache {
	return &RedisTableCache{
		client: client,
	}
}

func (r *RedisTableCache) GetCachedTable(key string) (*domain.CurrentTableSession, error) {
	data, err := r.client.Get(key)
	if err != nil {
		return nil, err
	}

	var table domain.CurrentTableSession
	err = json.Unmarshal([]byte(data), &table)
	if err != nil {
		return nil, err
	}

	return &table, nil
}

func (r *RedisTableCache) SetCachedTable(key string, table *domain.CurrentTableSession, ttl time.Duration) error {
	data, err := json.Marshal(table)
	if err != nil {
		return err
	}

	return r.client.Set(key, string(data), ttl)
}

func (r *RedisTableCache) DeleteCachedTable(key string) error {
	return r.client.Del(key)
}
