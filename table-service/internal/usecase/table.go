package usecase

import (
	"context"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"food-story/shared/redis"
	"food-story/table-service/internal/domain"
	"github.com/google/uuid"
	"time"
)

func (i *TableImplement) ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError) {
	return i.repository.ListTableStatus(ctx)
}

func (i *TableImplement) CreateTable(ctx context.Context, payload domain.Table) (result int64, customError *exceptions.CustomError) {
	return i.repository.CreateTable(ctx, payload)
}

func (i *TableImplement) UpdateTable(ctx context.Context, payload domain.Table) (customError *exceptions.CustomError) {

	customError = i.repository.IsTableExists(ctx, payload.ID)
	if customError != nil {
		return
	}

	return i.repository.UpdateTables(ctx, payload)
}

func (i *TableImplement) UpdateTableStatus(ctx context.Context, payload domain.TableStatus) (customError *exceptions.CustomError) {

	customError = i.repository.IsTableExists(ctx, payload.ID)
	if customError != nil {
		return
	}

	customError = i.repository.UpdateTablesStatus(ctx, payload)
	if customError != nil {
		return
	}

	return nil
}

func (i *TableImplement) SearchTableByFilters(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError) {
	return i.repository.SearchTables(ctx, payload)
}

func (i *TableImplement) QuickSearchAvailableTable(ctx context.Context, payload domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
	return i.repository.QuickSearchTables(ctx, payload)
}

func (i *TableImplement) CreateTableSession(ctx context.Context, payload domain.TableSession) (result string, customError *exceptions.CustomError) {

	customError = i.repository.IsTableAvailableOrReserved(ctx, payload.TableID)
	if customError != nil {
		return
	}

	sessionID := uuid.New()
	expiry := time.Now().Add(1 * time.Hour)
	sessionIDEncrypt, err := utils.EncryptSession(utils.SessionData{
		SessionID: sessionID.String(),
		Expiry:    expiry,
	}, []byte(i.config.SecretKey))
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: fmt.Errorf("failed to create table session: %w", err),
		}
	}

	tableNumber, customError := i.repository.GetTableNumber(ctx, payload.TableID)
	if customError != nil {
		return "", customError
	}

	key := redis.KeyTable + sessionID.String()
	err = i.cache.SetCachedTable(key, &domain.CurrentTableSession{
		SessionID:   sessionID,
		TableID:     payload.TableID,
		TableNumber: tableNumber,
		Status:      "active",
		StartedAt:   time.Now(),
		OrderID:     nil,
	}, 1*time.Hour)
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: fmt.Errorf("failed to create table session: %w", err),
		}
	}

	customError = i.repository.CreateTableSession(ctx, payload, sessionID, expiry)
	if customError != nil {
		_ = i.cache.DeleteCachedTable(key)
		return
	}

	return i.config.FrontendURL + "?s=" + sessionIDEncrypt, nil
}

func (i *TableImplement) GetCurrentSession(sessionIDEncrypt string) (*domain.CurrentTableSession, *exceptions.CustomError) {
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

	key := "table:" + sessionID
	cachedTable, err := i.cache.GetCachedTable(key)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: fmt.Errorf("failed to get current session: %w", err),
		}
	}

	if cachedTable == nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: fmt.Errorf("session not found"),
		}
	}

	return cachedTable, nil
}
