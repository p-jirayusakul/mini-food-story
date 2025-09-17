package cache

import (
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/shared/redis"

	"github.com/google/uuid"
)

type RedisTableCacheInterface interface {
	IsCachedTableExist(sessionID uuid.UUID) *exceptions.CustomError
}

type RedisTableCache struct {
	client redis.RedisInterface
}

func NewRedisTableCache(client *redis.RedisClient) *RedisTableCache {
	return &RedisTableCache{
		client: client,
	}
}

func (r *RedisTableCache) IsCachedTableExist(sessionID uuid.UUID) *exceptions.CustomError {
	_, err := r.client.Get(redis.KeyTable + sessionID.String())
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return &exceptions.CustomError{
				Status: exceptions.ERRAUTHORIZED,
				Errors: exceptions.ErrSessionExpired,
			}
		}
		return &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: fmt.Errorf("cache session: %w", err),
		}
	}

	return nil
}
