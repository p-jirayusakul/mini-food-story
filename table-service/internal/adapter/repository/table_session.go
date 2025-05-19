package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"food-story/table-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
	"time"
)

func (i *Implement) CreateTableSession(ctx context.Context, payload domain.TableSession, sessionID uuid.UUID, expiry time.Time) *exceptions.CustomError {

	err := i.repository.TXCreateTableSession(ctx, database.CreateTableSessionParams{
		ID:             i.snowflakeID.Generate(),
		TableID:        payload.TableID,
		NumberOfPeople: payload.NumberOfPeople,
		SessionID:      utils.UUIDToPgUUID(sessionID),
		ExpireAt:       pgtype.Timestamptz{Time: expiry, Valid: true},
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to create table session: %w", err),
		}
	}

	return nil
}

func (i *Implement) GetTableSession(ctx context.Context, sessionID uuid.UUID) (*domain.CurrentTableSession, *exceptions.CustomError) {

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

	result := domain.CurrentTableSession{
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
