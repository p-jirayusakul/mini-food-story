package usecase

import (
	"context"
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

	tableNumber, customError := i.repository.GetTableNumber(ctx, payload.TableID)
	if customError != nil {
		return "", customError
	}

	key := redis.KeyTable + sessionID.String()
	customError = i.cache.SetCachedTable(redis.KeyTable+sessionID.String(), &domain.CurrentTableSession{
		SessionID:   sessionID,
		TableID:     payload.TableID,
		TableNumber: tableNumber,
		Status:      SessionStatusActive,
		StartedAt:   time.Now(),
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
