package usecase

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"
	"food-story/shared/redis"
	"food-story/table-service/internal/domain"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const (
	SessionStatusActive = "active"
)

func (i *Implement) CreateTableSession(ctx context.Context, payload domain.TableSession) (result string, customError *exceptions.CustomError) {

	if payload.TableID <= 0 || payload.NumberOfPeople <= 0 {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: errors.New("invalid table session payload"),
		}
	}

	customError = i.repository.IsTableAvailableOrReserved(ctx, payload.TableID)
	if customError != nil {
		return "", customError
	}

	sessionID, sessionExpiry := i.generateSessionDetails()
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
	customError = i.cache.SetCachedTable(key, &shareModel.CurrentTableSession{
		SessionID:   sessionID,
		TableID:     payload.TableID,
		TableNumber: tableNumber,
		Status:      SessionStatusActive,
		StartedAt:   startedAt,
		OrderID:     nil,
	}, i.config.TableSessionDuration)
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

func (i *Implement) generateSessionDetails() (uuid.UUID, time.Time) {
	return uuid.New(), time.Now().Add(i.config.TableSessionDuration)
}

func encryptSessionID(sessionID uuid.UUID, expiry time.Time, key string) (string, *exceptions.CustomError) {
	result, err := utils.EncryptSession(utils.SessionData{
		SessionID: sessionID.String(),
		Expiry:    expiry,
	}, []byte(key))

	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
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
