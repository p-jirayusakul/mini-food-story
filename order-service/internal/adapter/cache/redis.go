package cache

import (
	"encoding/json"
	"errors"
	"food-story/order-service/internal/domain"
	"food-story/shared/redis"
	"github.com/google/uuid"
	"strconv"
)

type RedisTableCacheInterface interface {
	GetCachedTable(sessionID uuid.UUID) (*domain.CurrentTableSession, error)
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

func (r *RedisTableCache) UpdateOrderID(sessionID uuid.UUID, orderID int64) error {
	key := redis.KeyTable + sessionID.String()

	// ดึง TTL เดิม
	oldTTL, err := r.client.TTL(key)
	if err != nil {
		return err
	}
	if oldTTL <= 0 {
		return errors.New("session is expired")
	}

	table, err := r.GetCachedTable(sessionID)
	if err != nil {
		return err
	}

	orderIDStr := strconv.FormatInt(orderID, 10)
	table.OrderID = &orderIDStr

	data, err := json.Marshal(table)
	if err != nil {
		return err
	}

	return r.client.Set(key, string(data), oldTTL)
}
