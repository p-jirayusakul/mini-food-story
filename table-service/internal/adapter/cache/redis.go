package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	shareModel "food-story/shared/model"
	"food-story/shared/redis"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type RedisTableCacheInterface interface {
	GetCachedTable(sessionID uuid.UUID) (*shareModel.CurrentTableSession, *exceptions.CustomError)
	SetCachedTable(key string, table *shareModel.CurrentTableSession, ttl time.Duration) *exceptions.CustomError
	DeleteCachedTable(key string) *exceptions.CustomError
	IsCachedTableExist(sessionID uuid.UUID) *exceptions.CustomError
	SetCachedTableNumber(key string, tableNumber int32, ttl time.Duration) *exceptions.CustomError
	GetCachedTableNumber(key string) (int32, *exceptions.CustomError)
	GetTTL(sessionID uuid.UUID) (time.Duration, *exceptions.CustomError)
	ExtensionTTL(sessionID uuid.UUID, newTTL time.Duration) *exceptions.CustomError
}

type RedisTableCache struct {
	client redis.RedisInterface
}

func NewRedisTableCache(client *redis.RedisClient) *RedisTableCache {
	return &RedisTableCache{
		client: client,
	}
}

func (r *RedisTableCache) GetCachedTable(sessionID uuid.UUID) (*shareModel.CurrentTableSession, *exceptions.CustomError) {
	data, err := r.client.Get(redis.KeyTable + sessionID.String())
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrSessionNotFound,
			}
		}
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: err,
		}
	}

	var table shareModel.CurrentTableSession
	err = json.Unmarshal([]byte(data), &table)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: err,
		}
	}

	return &table, nil
}

func (r *RedisTableCache) GetTTL(sessionID uuid.UUID) (time.Duration, *exceptions.CustomError) {
	key := redis.KeyTable + sessionID.String()
	ttl, err := r.client.TTL(key)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: err,
		}
	}

	if ttl == -1 {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: errors.New("key has no expiration (persist)"),
		}
	} else if ttl == -2 {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: exceptions.ErrSessionNotFound,
		}
	}

	return ttl, nil
}

func (r *RedisTableCache) ExtensionTTL(sessionID uuid.UUID, newTTL time.Duration) *exceptions.CustomError {
	key := redis.KeyTable + sessionID.String()
	sessionDetail, sessionErr := r.GetCachedTable(sessionID)
	if sessionErr != nil {
		return sessionErr
	}

	delErr := r.DeleteCachedTable(key)
	if delErr != nil {
		return delErr
	}

	setCacheErr := r.SetCachedTable(key, sessionDetail, newTTL)
	if setCacheErr != nil {
		return setCacheErr
	}

	return nil
}

func (r *RedisTableCache) SetCachedTable(key string, table *shareModel.CurrentTableSession, ttl time.Duration) *exceptions.CustomError {
	data, err := json.Marshal(table)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: fmt.Errorf("pare json session: %w", err),
		}
	}

	err = r.client.Set(key, string(data), ttl)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: fmt.Errorf("set cache session: %w", err),
		}
	}

	return nil
}

func (r *RedisTableCache) SetCachedTableNumber(key string, tableNumber int32, ttl time.Duration) *exceptions.CustomError {

	err := r.client.Set(key, strconv.Itoa(int(tableNumber)), ttl)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: fmt.Errorf("set cache table number: %w", err),
		}
	}

	return nil
}

func (r *RedisTableCache) GetCachedTableNumber(key string) (int32, *exceptions.CustomError) {
	data, err := r.client.Get(key)
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return 0, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrRedisKeyNotFound,
			}
		}

		return 0, &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: fmt.Errorf("get cache session: %w", err),
		}
	}

	parsedValue, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: fmt.Errorf("parese string to int32: %w", err),
		}
	}

	return int32(parsedValue), nil
}

func (r *RedisTableCache) DeleteCachedTable(key string) *exceptions.CustomError {
	err := r.client.Del(key)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: fmt.Errorf("delete cache session: %w", err),
		}
	}
	return nil
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
