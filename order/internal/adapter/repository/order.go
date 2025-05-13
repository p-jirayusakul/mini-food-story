package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/order/internal/domain"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (i *Implement) CreateOrder(ctx context.Context, order domain.Order) (result int64, customError *exceptions.CustomError) {

	var sessionByte [16]byte = order.SessionID
	id, err := i.repository.CreateOrder(ctx, database.CreateOrderParams{
		ID:        i.snowflakeID.Generate(),
		SessionID: pgtype.UUID{Bytes: sessionByte, Valid: true},
		TableID:   order.TableID,
	})
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to create order: %w", err),
		}
	}

	return id, nil
}

func (i *Implement) GetOrderByID(ctx context.Context, id int64) (result *domain.Order, customError *exceptions.CustomError) {
	data, err := i.repository.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: fmt.Errorf("order not found"),
			}
		}
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check order exists: %w", err),
		}
	}

	sessionID, err := uuid.Parse(data.SessionID.String())
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: fmt.Errorf("failed to get order: %w", err),
		}
	}

	return &domain.Order{
		ID:           data.ID,
		SessionID:    sessionID,
		TableID:      data.TableID,
		StatusID:     data.StatusID,
		StatusName:   data.StatusName,
		StatusNameEN: data.StatusNameEN,
	}, nil
}

func (i *Implement) UpdateOrderStatus(ctx context.Context, payload domain.OrderStatus) (customError *exceptions.CustomError) {
	customError = i.IsOrderExist(ctx, payload.ID)
	if customError != nil {
		return
	}

	customError = i.IsOrderStatus(ctx, payload.StatusCode)
	if customError != nil {
		return
	}

	err := i.repository.UpdateOrderStatus(ctx, database.UpdateOrderStatusParams{
		StatusCode: payload.StatusCode,
		ID:         payload.ID,
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update order status: %w", err),
		}
	}

	return
}

func (i *Implement) IsOrderExist(ctx context.Context, id int64) (customError *exceptions.CustomError) {
	isExist, err := i.repository.IsOrderExist(ctx, id)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check order exists: %w", err),
		}
	}

	if !isExist {
		return &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: fmt.Errorf("order not found"),
		}
	}

	return nil
}
