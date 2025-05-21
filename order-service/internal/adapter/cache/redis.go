package cache

import (
	"encoding/json"
	"errors"
	"food-story/pkg/exceptions"
	shareModel "food-story/shared/model"
	"food-story/shared/redis"
	"github.com/google/uuid"
	"strconv"
)

type RedisTableCacheInterface interface {
	GetCachedTable(sessionID uuid.UUID) (*shareModel.CurrentTableSession, *exceptions.CustomError)
	DeleteCachedTable(sessionID uuid.UUID) *exceptions.CustomError
	IsCachedTableExist(sessionID uuid.UUID) *exceptions.CustomError
	UpdateOrderID(sessionID uuid.UUID, orderID int64) *exceptions.CustomError
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

func (r *RedisTableCache) IsCachedTableExist(sessionID uuid.UUID) *exceptions.CustomError {
	_, err := r.client.Get(redis.KeyTable + sessionID.String())
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrSessionNotFound,
			}
		}

		return &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: err,
		}
	}

	return nil
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

func (r *RedisTableCache) UpdateOrderID(sessionID uuid.UUID, orderID int64) *exceptions.CustomError {
	key := redis.KeyTable + sessionID.String()

	// ดึง TTL เดิม
	oldTTL, err := r.client.TTL(key)
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrSessionNotFound,
			}
		}

		return &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: err,
		}
	}
	if oldTTL <= 0 {
		return &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: errors.New("session is expired"),
		}
	}

	table, getTableErr := r.GetCachedTable(sessionID)
	if getTableErr != nil {
		return getTableErr
	}

	orderIDStr := strconv.FormatInt(orderID, 10)
	table.OrderID = &orderIDStr

	data, err := json.Marshal(table)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: err,
		}
	}

	err = r.client.Set(key, string(data), oldTTL)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRCACHE,
			Errors: err,
		}
	}

	return nil
}
