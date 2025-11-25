package cache

import (
	"encoding/json"
	"errors"
	"food-story/pkg/exceptions"
	shareModel "food-story/shared/model"
	"food-story/shared/redis"
	"strconv"

	"github.com/google/uuid"
)

type RedisTableCacheInterface interface {
	GetCachedTable(sessionID uuid.UUID) (*shareModel.CurrentTableSession, error)
	DeleteCachedTable(sessionID uuid.UUID) error
	IsCachedTableExist(sessionID uuid.UUID) error
	UpdateOrderID(sessionID uuid.UUID, orderID int64) error
}

type RedisTableCache struct {
	client redis.RedisInterface
}

func NewRedisTableCache(client *redis.RedisClient) *RedisTableCache {
	return &RedisTableCache{
		client: client,
	}
}

func (r *RedisTableCache) GetCachedTable(sessionID uuid.UUID) (*shareModel.CurrentTableSession, error) {
	data, err := r.client.Get(redis.KeyTable + sessionID.String())
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return nil, exceptions.ErrRedisKeyNotFoundException
		}
		return nil, exceptions.Errorf(exceptions.CodeRedis, "failed to get cached table", err)
	}

	var table shareModel.CurrentTableSession
	err = json.Unmarshal([]byte(data), &table)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRedis, "failed to unmarshal cached table", err)
	}

	return &table, nil
}

func (r *RedisTableCache) IsCachedTableExist(sessionID uuid.UUID) error {
	_, err := r.client.Get(redis.KeyTable + sessionID.String())
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return exceptions.ErrRedisKeyNotFoundException
		}

		return exceptions.Errorf(exceptions.CodeRedis, "failed to get cached table", err)
	}

	return nil
}

func (r *RedisTableCache) DeleteCachedTable(sessionID uuid.UUID) error {
	err := r.client.Del(redis.KeyTable + sessionID.String())
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRedis, "failed to delete cached table", err)
	}
	return nil
}

func (r *RedisTableCache) UpdateOrderID(sessionID uuid.UUID, orderID int64) error {
	key := redis.KeyTable + sessionID.String()

	// ดึง TTL เดิม
	oldTTL, err := r.client.TTL(key)
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return exceptions.ErrRedisKeyNotFoundException
		}

		return exceptions.Errorf(exceptions.CodeRedis, "failed to get ttl table", err)
	}
	if oldTTL <= 0 {
		return exceptions.Error(exceptions.CodeRedis, "session is expired")
	}

	table, getTableErr := r.GetCachedTable(sessionID)
	if getTableErr != nil {
		return getTableErr
	}

	orderIDStr := strconv.FormatInt(orderID, 10)
	table.OrderID = &orderIDStr

	data, err := json.Marshal(table)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeSystem, "failed to marshal cached table", err)
	}

	err = r.client.Set(key, string(data), oldTTL)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRedis, "failed to set cached table", err)

	}

	return nil
}
