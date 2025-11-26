package usecase

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"food-story/shared/redis"
	"food-story/table-service/internal/domain"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const (
	_sessionStatusActive    = "active"
	_errInvalidTableSession = "invalid table session payload"
)

func (i *Implement) ListSessionExtensionReason(ctx context.Context) (result []*domain.ListSessionExtensionReason, err error) {
	return i.repository.ListSessionExtensionReason(ctx)
}

func (i *Implement) CreateTableSession(ctx context.Context, payload domain.TableSession) (result string, err error) {

	if payload.TableID <= 0 || payload.NumberOfPeople <= 0 {
		return "", exceptions.Error(exceptions.CodeBusiness, _errInvalidTableSession)
	}

	err = i.repository.IsTableAvailableOrReserved(ctx, payload.TableID)
	if err != nil {
		return "", err
	}

	sessionID, sessionExpiry := i.generateSessionDetails()
	encryptedSessionID, err := encryptSessionID(sessionID, sessionExpiry, i.config.SecretKey)
	if err != nil {
		return "", err
	}

	tableNumber, err := i.getTableNumberFromCache(ctx, payload.TableID)
	if err != nil {
		return "", err
	}

	startedAt, err := i.repository.GetCurrentDateTime(ctx)
	if err != nil {
		return "", err
	}

	key := redis.KeyTable + sessionID.String()
	err = i.cache.SetCachedTable(key, &domain.CurrentTableSession{
		SessionID:   sessionID,
		TableID:     payload.TableID,
		TableNumber: tableNumber,
		Status:      _sessionStatusActive,
		StartedAt:   startedAt,
		OrderID:     nil,
	}, i.config.TableSessionDuration)
	if err != nil {
		return "", err
	}

	err = i.repository.CreateTableSession(ctx, payload, sessionID, sessionExpiry)
	if err != nil {
		cacheErr := i.cache.DeleteCachedTable(key)
		if cacheErr != nil {
			slog.Error("failed to delete cache table session: ", "error", cacheErr)
		}
		return "", err
	}

	return i.config.FrontendURL + "?s=" + encryptedSessionID, nil
}

func (i *Implement) generateSessionDetails() (uuid.UUID, time.Time) {
	return uuid.New(), time.Now().Add(i.config.TableSessionDuration)
}

func encryptSessionID(sessionID uuid.UUID, expiry time.Time, key string) (string, error) {
	result, err := utils.EncryptSession(utils.SessionData{
		SessionID: sessionID.String(),
		Expiry:    expiry,
	}, []byte(key))

	if err != nil {
		return "", exceptions.Errorf(exceptions.CodeSystem, "failed to encrypt session ID", err)
	}

	return result, nil
}

func (i *Implement) getTableNumberFromCache(ctx context.Context, tableID int64) (int32, error) {
	keyTableNumber := fmt.Sprintf("%s:%d", redis.KeyTable, tableID)
	tableNumber, err := i.cache.GetCachedTableNumber(keyTableNumber)
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFoundException) {
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
			return 0, err
		}
	}

	return tableNumber, nil
}

func (i *Implement) SessionExtension(ctx context.Context, payload domain.SessionExtension) error {

	requestedMinutes, err := i.repository.GetDurationMinutesByProductID(ctx, payload.ProductID)
	if err != nil {
		return err
	}

	sessionID, err := i.repository.GetSessionIDByTableID(ctx, payload.TableID)
	if err != nil {
		return err
	}

	oldTTL, err := i.cache.GetTTL(sessionID)
	if err != nil {
		return err
	}

	newTTLMinutes := oldTTL + time.Duration(requestedMinutes)*time.Minute
	err = i.cache.ExtensionTTL(sessionID, newTTLMinutes)
	if err != nil {
		return err
	}

	err = i.repository.SessionExtension(ctx, payload, requestedMinutes)
	if err != nil {
		extensionTTLErr := i.cache.ExtensionTTL(sessionID, oldTTL)
		if extensionTTLErr != nil {
			return extensionTTLErr
		}
		return err
	}

	return nil
}
