package cache

import (
	"errors"
	"food-story/pkg/exceptions"
	"food-story/shared/redis"

	"github.com/google/uuid"
)

type RedisTableCacheInterface interface {
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

func (r *RedisTableCache) IsCachedTableExist(sessionID uuid.UUID) error {
	_, err := r.client.Get(redis.KeyTable + sessionID.String())
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return exceptions.Error(exceptions.CodeUnauthorized, exceptions.ErrSessionExpired.Error())
		}
		return exceptions.Errorf(exceptions.CodeRedis, "cache session", err)
	}

	return nil
}
