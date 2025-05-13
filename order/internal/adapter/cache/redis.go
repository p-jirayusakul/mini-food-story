package cache

import (
	"encoding/json"
	"food-story/order/internal/domain"
	"food-story/shared/redis"
	"github.com/google/uuid"
)

type RedisTableCacheInterface interface {
	GetCachedTable(sessionID uuid.UUID) (*domain.CurrentTableSession, error)
	DeleteCachedTable(sessionID uuid.UUID) error
	IsCachedTableExist(sessionID uuid.UUID) error
}

type RedisTableCache struct {
	client redis.RedisInterface
}

func NewRedisTableCache(client *redis.RedisClient) *RedisTableCache {
	return &RedisTableCache{
		client: client,
	}
}

func (r *RedisTableCache) GetCachedTable(sessionID uuid.UUID) (*domain.CurrentTableSession, error) {
	data, err := r.client.Get(redis.KeyTable + sessionID.String())
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

func (r *RedisTableCache) IsCachedTableExist(sessionID uuid.UUID) error {
	_, err := r.client.Get(redis.KeyTable + sessionID.String())
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisTableCache) DeleteCachedTable(sessionID uuid.UUID) error {
	return r.client.Del(redis.KeyTable + sessionID.String())
}
