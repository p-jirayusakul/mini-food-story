package repository

import (
	"context"
	"errors"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"food-story/table-service/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (i *Implement) ListSessionExtensionReason(ctx context.Context) (result []*domain.ListSessionExtensionReason, err error) {
	const _errMessage = "failed to fetch session extension reason"

	data, err := i.repository.ListSessionExtensionReason(ctx)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, _errMessage, err)
	}

	if data == nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, _errMessage, err)
	}

	result = make([]*domain.ListSessionExtensionReason, len(data))
	for index, v := range data {
		result[index] = &domain.ListSessionExtensionReason{
			ID:       v.ID,
			Code:     v.Code,
			Name:     v.Name,
			NameEN:   v.NameEN,
			Category: utils.PgTextToStringPtr(v.Category),
			ModeCode: utils.PgTextToStringPtr(v.ModeCode),
		}
	}

	return result, nil
}

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
			return uuid.UUID{}, exceptions.ErrorSessionNotFound()
		}
		return uuid.UUID{}, exceptions.Errorf(exceptions.CodeRepository, "failed to get session id by table id", err)
	}

	v, err := utils.PareStringToUUID(sessionID.String())
	if err != nil {
		return uuid.UUID{}, exceptions.Errorf(exceptions.CodeSystem, "failed to parse session id", err)
	}

	return v, nil
}

func (i *Implement) SessionExtension(ctx context.Context, payload domain.SessionExtension, requestedMinutes int32) error {
	tableExp, err := i.repository.GetExpiresAtByTableID(ctx, payload.TableID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return exceptions.ErrorSessionNotFound()
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
		return exceptions.ErrorSessionNotFound()
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
		NewOrderItemsID:        i.snowflakeID.Generate(),
	})
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to create session extension", err)
	}
	return nil
}

func (i *Implement) GetOrderIDBySessionID(ctx context.Context, sessionID uuid.UUID) (orderID int64, err error) {
	data, err := i.repository.GetOrderIDBySessionID(ctx, utils.UUIDToPgUUID(sessionID))
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, exceptions.ErrorSessionNotFound()
		}
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch order id by session id", err)
	}
	return data, nil
}

func (i *Implement) GetDurationMinutesByProductID(ctx context.Context, productID int64) (durationMinutes int32, err error) {
	data, err := i.repository.GetDurationMinutesByProductID(ctx, productID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, exceptions.ErrorIDNotFound(exceptions.CodeProductNotFound, productID)
		}
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch duration minutes by product id", err)
	}
	return data, nil
}
