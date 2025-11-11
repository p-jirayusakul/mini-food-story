package repository

import (
	"context"
	"errors"
	"fmt"
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

func (i *Implement) CreateTableSession(ctx context.Context, payload domain.TableSession, sessionID uuid.UUID, expiry time.Time) *exceptions.CustomError {

	err := i.repository.TXCreateTableSession(ctx, database.CreateTableSessionParams{
		ID:             i.snowflakeID.Generate(),
		TableID:        payload.TableID,
		NumberOfPeople: payload.NumberOfPeople,
		SessionID:      utils.UUIDToPgUUID(sessionID),
		ExpiresAt:      pgtype.Timestamptz{Time: expiry, Valid: true},
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to create table session: %w", err),
		}
	}

	return nil
}

func (i *Implement) GetTableSession(ctx context.Context, sessionID uuid.UUID) (*shareModel.CurrentTableSession, *exceptions.CustomError) {

	data, err := i.repository.GetTableSession(ctx, utils.UUIDToPgUUID(sessionID))
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: errors.New("table session not found"),
			}
		}
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get table session: %w", err),
		}
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

func (i *Implement) GetCurrentDateTime(ctx context.Context) (time.Time, *exceptions.CustomError) {
	currentTime, err := i.repository.GetTimeNow(ctx)
	if err != nil {
		return time.Time{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get current time: %w", err),
		}
	}

	return currentTime.Time, nil
}

func (i *Implement) GetSessionIDByTableID(ctx context.Context, tableID int64) (uuid.UUID, *exceptions.CustomError) {
	sessionID, err := i.repository.GetSessionIDByTableID(ctx, tableID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return uuid.UUID{}, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: errors.New("table session id not found"),
			}
		}
		return uuid.UUID{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get session id: %w", err),
		}
	}

	v, err := utils.PareStringToUUID(sessionID.String())
	if err != nil {
		return uuid.UUID{}, &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: fmt.Errorf("failed to parse session id: %w", err),
		}
	}

	return v, nil
}

func (i *Implement) SessionExtension(ctx context.Context, payload domain.SessionExtension) *exceptions.CustomError {
	tableExp, err := i.repository.GetExpiresAtByTableID(ctx, payload.TableID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: errors.New("table session not found"),
			}
		}
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get expires at by table id: %w", err),
		}
	}

	if tableExp.ExtendTotalMinutes >= tableExp.MaxExtendMinutes {
		return &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: errors.New("table extension is not allowed"),
		}
	} else if int32(payload.RequestedMinutes) > tableExp.MaxExtendMinutes {
		return &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: errors.New("requested minutes is not allowed"),
		}
	}

	expireAt, err := utils.PgTimestampToTime(tableExp.ExpiresAt)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: fmt.Errorf("failed to parse expires at: %w", err),
		}
	}

	totalExpiresAt := expireAt.Add(time.Duration(payload.RequestedMinutes) * time.Minute)

	reasonResult, err := i.repository.GetSessionExtensionModeByReasonCode(ctx, payload.ReasonCode)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get session extension mode by reason code: %w", err),
		}
	}

	if reasonResult.SessionExtensionModeID.Int64 == 0 {
		return &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: errors.New("session extension mode not found"),
		}
	}

	createSessionExtensionParams := database.CreateSessionExtensionParams{
		ID:               i.snowflakeID.Generate(),
		RequestedMinutes: int32(payload.RequestedMinutes),
		ReasonID:         utils.Int64ToPgInt8(reasonResult.ID),
		ModeID:           utils.Int64ToPgInt8(reasonResult.SessionExtensionModeID.Int64),
	}

	err = i.repository.TXSessionsExtension(ctx, database.TXSessionsExtensionParams{
		TableID:                payload.TableID,
		RequestedMinutes:       payload.RequestedMinutes,
		ReasonCode:             payload.ReasonCode,
		ExpiresAt:              totalExpiresAt,
		CreateSessionExtension: createSessionExtensionParams,
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to session extension: %w", err),
		}
	}
	return nil
}
