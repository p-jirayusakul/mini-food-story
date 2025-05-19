package usecase

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"food-story/shared/redis"
	"food-story/table-service/internal/domain"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

const (
	TableSessionDuration = 1 * time.Hour
	SessionStatusActive  = "active"
)

func (i *Implement) CreateTableSession(ctx context.Context, payload domain.TableSession) (result string, customError *exceptions.CustomError) {

	if customError = i.repository.IsTableAvailableOrReserved(ctx, payload.TableID); customError != nil {
		return "", customError
	}

	sessionID, sessionExpiry := generateSessionDetails()
	encryptedSessionID, customError := encryptSessionID(sessionID, sessionExpiry, i.config.SecretKey)
	if customError != nil {
		return "", customError
	}

	tableNumber, customError := i.getTableNumberFromCache(ctx, payload.TableID)
	if customError != nil {
		return "", customError
	}

	startedAt, customError := i.repository.GetCurrentDateTime(ctx)
	if customError != nil {
		return "", customError
	}

	key := redis.KeyTable + sessionID.String()
	customError = i.cache.SetCachedTable(key, &domain.CurrentTableSession{
		SessionID:   sessionID,
		TableID:     payload.TableID,
		TableNumber: tableNumber,
		Status:      SessionStatusActive,
		StartedAt:   startedAt,
		OrderID:     nil,
	}, TableSessionDuration)
	if customError != nil {
		return "", customError
	}

	customError = i.repository.CreateTableSession(ctx, payload, sessionID, sessionExpiry)
	if customError != nil {
		cacheErr := i.cache.DeleteCachedTable(key)
		if cacheErr != nil {
			slog.Error("failed to delete cache table session: ", "error", cacheErr)
		}
		return "", customError
	}

	return i.config.FrontendURL + "?s=" + encryptedSessionID, nil
}

func (i *Implement) GetCurrentSession(sessionIDEncrypt string) (*domain.CurrentTableSession, *exceptions.CustomError) {
	sessionIDDecrypt, err := utils.DecryptSession(sessionIDEncrypt, []byte(i.config.SecretKey))
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: fmt.Errorf("failed to get current session: %w", err),
		}
	}

	sessionID := sessionIDDecrypt.SessionID
	expiry := sessionIDDecrypt.Expiry
	if expiry.Before(time.Now()) {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: fmt.Errorf("session expired"),
		}
	}

	cachedTable, customError := i.cache.GetCachedTable(redis.KeyTable + sessionID)
	if customError != nil {
		return nil, customError
	}

	if cachedTable == nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: fmt.Errorf("session not found"),
		}
	}

	return cachedTable, nil
}

func generateSessionDetails() (uuid.UUID, time.Time) {
	return uuid.New(), time.Now().Add(TableSessionDuration)
}

func encryptSessionID(sessionID uuid.UUID, expiry time.Time, key string) (string, *exceptions.CustomError) {
	result, err := utils.EncryptSession(utils.SessionData{
		SessionID: sessionID.String(),
		Expiry:    expiry,
	}, []byte(key))

	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: fmt.Errorf("failed to encrypt session ID: %w", err),
		}
	}

	return result, nil
}

func (i *Implement) getTableNumberFromCache(ctx context.Context, tableID int64) (int32, *exceptions.CustomError) {
	keyTableNumber := fmt.Sprintf("%s:%d", redis.KeyTable, tableID)
	tableNumber, customError := i.cache.GetCachedTableNumber(keyTableNumber)
	if customError != nil {
		if errors.Is(customError.Errors, exceptions.ErrRedisKeyNotFound) {
			tableNumberDB, getTableNumberErr := i.repository.GetTableNumber(ctx, tableID)
			if getTableNumberErr != nil {
				return 0, getTableNumberErr
			}
			setTableNumberErr := i.cache.SetCachedTableNumber(keyTableNumber, tableNumberDB, 24*time.Hour)
			if setTableNumberErr != nil {
				return 0, setTableNumberErr
			}
			tableNumber = tableNumberDB
		} else {
			return 0, customError
		}
	}

	if tableNumber == 0 {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: errors.New("table not found"),
		}
	}

	return tableNumber, nil
}
