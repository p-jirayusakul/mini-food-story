package cache

import (
	"encoding/json"
	"errors"
	"food-story/pkg/exceptions"
	"food-story/shared/redis"
	"food-story/table-service/internal/domain"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type RedisTableCacheInterface interface {
	GetCachedTableStatus() ([]*domain.Status, error)
	SetCachedTableStatus(payload []*domain.Status) error
	GetCachedTable(sessionID uuid.UUID) (*domain.CurrentTableSession, error)
	SetCachedTable(key string, table *domain.CurrentTableSession, ttl time.Duration) error
	DeleteCachedTable(key string) error
	IsCachedTableExist(sessionID uuid.UUID) error
	SetCachedTableNumber(key string, tableNumber int32, ttl time.Duration) error
	GetCachedTableNumber(key string) (int32, error)
	GetTTL(sessionID uuid.UUID) (time.Duration, error)
	ExtensionTTL(sessionID uuid.UUID, newTTL time.Duration) error
}

type RedisTableCache struct {
	client redis.RedisInterface
}

func NewRedisTableCache(client *redis.RedisClient) *RedisTableCache {
	return &RedisTableCache{
		client: client,
	}
}

func (r *RedisTableCache) GetCachedTableStatus() ([]*domain.Status, error) {
	data, err := r.client.Get(redis.KeyTableStatus)
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return nil, exceptions.ErrRedisKeyNotFoundException
		}
		return nil, exceptions.Errorf(exceptions.CodeRedis, "failed to get cache table status", err)
	}

	var status []*domain.Status
	err = json.Unmarshal([]byte(data), &status)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeSystem, "failed to unmarshal cache table status", err)
	}

	return status, nil
}

func (r *RedisTableCache) GetCachedTable(sessionID uuid.UUID) (*domain.CurrentTableSession, error) {
	data, err := r.client.Get(redis.KeyTable + sessionID.String())
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return nil, exceptions.ErrRedisKeyNotFoundException
		}
		return nil, exceptions.Errorf(exceptions.CodeRedis, "failed to get cache table", err)
	}

	var table domain.CurrentTableSession
	err = json.Unmarshal([]byte(data), &table)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeSystem, "failed to unmarshal cache table", err)
	}

	return &table, nil
}

func (r *RedisTableCache) SetCachedTableStatus(payload []*domain.Status) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeSystem, "failed to pare json session set cache table status", err)
	}

	err = r.client.Set(redis.KeyTableStatus, string(data), time.Hour*24)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRedis, "failed to set cache table status", err)
	}

	return nil
}

func (r *RedisTableCache) GetTTL(sessionID uuid.UUID) (time.Duration, error) {
	key := redis.KeyTable + sessionID.String()
	ttl, err := r.client.TTL(key)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRedis, "failed to get TTL", err)
	}

	switch ttl {
	case -1:
		return 0, exceptions.Error(exceptions.CodeRedis, "key has no expiration (persist)")
	case -2:
		return 0, exceptions.ErrorSessionNotFound()
	default:
		return ttl, nil
	}
}

func (r *RedisTableCache) ExtensionTTL(sessionID uuid.UUID, newTTL time.Duration) error {
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

func (r *RedisTableCache) SetCachedTable(key string, table *domain.CurrentTableSession, ttl time.Duration) error {
	data, err := json.Marshal(table)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeSystem, "failed to pare json session set cache table", err)
	}

	err = r.client.Set(key, string(data), ttl)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRedis, "failed to set cache table", err)
	}

	return nil
}

func (r *RedisTableCache) SetCachedTableNumber(key string, tableNumber int32, ttl time.Duration) error {

	err := r.client.Set(key, strconv.Itoa(int(tableNumber)), ttl)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRedis, "failed to set cache table number", err)
	}

	return nil
}

func (r *RedisTableCache) GetCachedTableNumber(key string) (int32, error) {
	data, err := r.client.Get(key)
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return 0, exceptions.ErrRedisKeyNotFoundException
		}

		return 0, exceptions.Errorf(exceptions.CodeRedis, "failed to get cache session", err)
	}

	parsedValue, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeSystem, "failed to parse string to int32", err)
	}

	return int32(parsedValue), nil
}

func (r *RedisTableCache) DeleteCachedTable(key string) error {
	err := r.client.Del(key)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRedis, "failed to delete cache table session", err)
	}
	return nil
}

func (r *RedisTableCache) IsCachedTableExist(sessionID uuid.UUID) error {
	_, err := r.client.Get(redis.KeyTable + sessionID.String())
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return exceptions.Error(exceptions.CodeUnauthorized, exceptions.ErrSessionExpired.Error())
		}
		return exceptions.Errorf(exceptions.CodeRedis, "failed to get cache table session", err)
	}

	return nil
}
