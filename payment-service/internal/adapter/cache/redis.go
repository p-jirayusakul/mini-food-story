package cache

import (
	"food-story/pkg/exceptions"
	"food-story/shared/redis"

	"github.com/google/uuid"
)

type RedisTableCacheInterface interface {
	DeleteCachedTable(sessionID uuid.UUID) *exceptions.CustomError
}

type RedisTableCache struct {
	client redis.RedisInterface
}

func NewRedisTableCache(client *redis.RedisClient) *RedisTableCache {
	return &RedisTableCache{
		client: client,
	}
}

func (r *RedisTableCache) DeleteCachedTable(sessionID uuid.UUID) *exceptions.CustomError {
	err := r.client.Del(redis.KeyTable + sessionID.String())
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: err,
		}
	}
	return nil
}
