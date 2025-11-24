package repository

import (
	"context"
	"errors"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	shareModel "food-story/shared/model"
	"food-story/table-service/internal/domain"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	_errTableSessionNotFound = exceptions.Error(exceptions.CodeNotFound, exceptions.ErrTableSessionNotFound.Error())
)

func (i *Implement) CreateTableSession(ctx context.Context, payload domain.TableSession, sessionID uuid.UUID, expiry time.Time) error {

	err := i.repository.TXCreateTableSession(ctx, database.CreateTableSessionParams{
		ID:             i.snowflakeID.Generate(),
		TableID:        payload.TableID,
		NumberOfPeople: payload.NumberOfPeople,
		SessionID:      utils.UUIDToPgUUID(sessionID),
		ExpiresAt:      pgtype.Timestamptz{Time: expiry, Valid: true},
	})
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to create table session", err)
	}

	return nil
}

func (i *Implement) GetTableSession(ctx context.Context, sessionID uuid.UUID) (*shareModel.CurrentTableSession, error) {

	data, err := i.repository.GetTableSession(ctx, utils.UUIDToPgUUID(sessionID))
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, _errTableSessionNotFound
		}
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to get table session", err)
	}

	result := shareModel.CurrentTableSession{
		SessionID:   sessionID,
		TableID:     data.TableID,
		TableNumber: data.TableNumber,
		Status:      string(data.Status.TableSessionStatus),
		StartedAt:   data.StartedAt.Time,
	}

	if data.OrderID.Valid {
		orderID := strconv.FormatInt(data.OrderID.Int64, 10)
		result.OrderID = &orderID
	}

	return &result, nil
}

func (i *Implement) GetCurrentDateTime(ctx context.Context) (time.Time, error) {
	currentTime, err := i.repository.GetTimeNow(ctx)
	if err != nil {
		return time.Time{}, exceptions.Errorf(exceptions.CodeRepository, "failed to get current time", err)
	}

	return currentTime.Time, nil
}

func (i *Implement) GetSessionIDByTableID(ctx context.Context, tableID int64) (uuid.UUID, error) {
	sessionID, err := i.repository.GetSessionIDByTableID(ctx, tableID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return uuid.UUID{}, _errTableSessionNotFound
		}
		return uuid.UUID{}, exceptions.Errorf(exceptions.CodeRepository, "failed to get session id by table id", err)
	}

	v, err := utils.PareStringToUUID(sessionID.String())
	if err != nil {
		return uuid.UUID{}, exceptions.Errorf(exceptions.CodeSystem, "failed to parse session id", err)
	}

	return v, nil
}

func (i *Implement) SessionExtension(ctx context.Context, payload domain.SessionExtension, requestedMinutes int32, orderID int64) error {
	tableExp, err := i.repository.GetExpiresAtByTableID(ctx, payload.TableID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return _errTableSessionNotFound
		}
		return exceptions.Errorf(exceptions.CodeRepository, "failed to get expires at by table id", err)
	}

	if tableExp.ExtendTotalMinutes >= tableExp.MaxExtendMinutes {
		return exceptions.Error(exceptions.CodeBusiness, "table extension is not allowed")
	} else if requestedMinutes > tableExp.MaxExtendMinutes {
		return exceptions.Error(exceptions.CodeBusiness, "requested minutes is not allowed")
	}

	expireAt, err := utils.PgTimestampToTime(tableExp.ExpiresAt)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeSystem, "failed to parse expires at", err)
	}

	totalExpiresAt := expireAt.Add(time.Duration(requestedMinutes) * time.Minute)

	reasonResult, err := i.repository.GetSessionExtensionModeByReasonCode(ctx, payload.ReasonCode)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to get session extension mode by reason code", err)
	}

	if reasonResult.SessionExtensionModeID.Int64 == 0 {
		return _errTableSessionNotFound
	}

	createSessionExtensionParams := database.CreateSessionExtensionParams{
		ID:               i.snowflakeID.Generate(),
		RequestedMinutes: requestedMinutes,
		ReasonID:         utils.Int64ToPgInt8(reasonResult.ID),
		ModeID:           utils.Int64ToPgInt8(reasonResult.SessionExtensionModeID.Int64),
	}

	err = i.repository.TXSessionsExtension(ctx, database.TXSessionsExtensionParams{
		TableID:                payload.TableID,
		RequestedMinutes:       int64(requestedMinutes),
		ReasonCode:             payload.ReasonCode,
		ExpiresAt:              totalExpiresAt,
		CreateSessionExtension: createSessionExtensionParams,
		ProductID:              payload.ProductID,
		OrderID:                orderID,
		NewOrderItemsID:        i.snowflakeID.Generate(),
	})
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to create session extension", err)
	}
	return nil
}
