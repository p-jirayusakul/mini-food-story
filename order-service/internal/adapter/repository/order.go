package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func (i *Implement) CreateOrder(ctx context.Context, payload domain.CreateOrder) (result int64, customError *exceptions.CustomError) {

	orderItemsPayload, customError := i.BuildPayloadOrderItems(ctx, payload.OrderItems)
	if customError != nil {
		return
	}

	orderNumber, customError := i.GetOrCreateOrderSequence(ctx)
	if customError != nil {
		return
	}

	var sessionByte [16]byte = payload.SessionID
	id, err := i.repository.TXCreateOrder(ctx, database.TXCreateOrderParams{
		CreateOrderItems: orderItemsPayload,
		CreateOrder: database.CreateOrderParams{
			ID:          i.snowflakeID.Generate(),
			OrderNumber: orderNumber,
			SessionID:   pgtype.UUID{Bytes: sessionByte, Valid: true},
			TableID:     payload.TableID,
		},
	})
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to create order: %w", err),
		}
	}

	return id, nil
}

func (i *Implement) GetOrCreateOrderSequence(ctx context.Context) (string, *exceptions.CustomError) {
	num, err := i.repository.GetOrCreateOrderSequence(ctx, pgtype.Date{
		Time:  time.Now(),
		Valid: true,
	})
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get or create order sequence: %w", err),
		}
	}

	return fmt.Sprintf("FS-%s-%04d", time.Now().Format("20060102"), num), nil
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

	return &domain.Order{
		ID:           data.ID,
		TableID:      data.TableID,
		TableNumber:  data.TableNumber,
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

func (i *Implement) IsOrderWithItemsExists(ctx context.Context, orderID, orderItemsID int64) (customError *exceptions.CustomError) {
	isExist, err := i.repository.IsOrderWithItemsExists(ctx, database.IsOrderWithItemsExistsParams{
		OrderID:      orderID,
		OrderItemsID: orderItemsID,
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check order item exists: %w", err),
		}
	}

	if !isExist {
		return &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: fmt.Errorf("order item not found"),
		}
	}

	return nil
}
