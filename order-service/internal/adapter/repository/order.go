package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (i *Implement) CreateOrder(ctx context.Context, order domain.CreateOrder) (orderID int64, customError *exceptions.CustomError) {

	orderItems, buildParamErr := i.buildPayloadOrderItems(ctx, order.OrderItems)
	if buildParamErr != nil {
		return 0, buildParamErr
	}

	orderNumber, sequenceError := i.getOrderSequence(ctx)
	if sequenceError != nil {
		return 0, sequenceError
	}

	orderID, err := i.repository.TXCreateOrder(ctx, database.TXCreateOrderParams{
		CreateOrderItems: orderItems,
		CreateOrder: database.CreateOrderParams{
			ID:          i.snowflakeID.Generate(),
			OrderNumber: orderNumber,
			SessionID:   utils.UUIDToPgUUID(order.SessionID),
			TableID:     order.TableID,
		},
	})
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to create order: %w", err),
		}
	}

	return orderID, nil
}

func (i *Implement) GetOrderByID(ctx context.Context, id int64) (result *domain.Order, customError *exceptions.CustomError) {
	order, err := i.repository.GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrOrderNotFound,
			}
		}
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check order exists: %w", err),
		}
	}

	if order == nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: exceptions.ErrOrderNotFound,
		}
	}

	return &domain.Order{
		ID:           order.ID,
		TableID:      order.TableID,
		TableNumber:  order.TableNumber,
		StatusID:     order.StatusID,
		StatusName:   order.StatusName,
		StatusNameEN: order.StatusNameEN,
		StatusCode:   order.StatusCode,
	}, nil
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
			Errors: exceptions.ErrOrderNotFound,
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
			Errors: exceptions.ErrOrderItemsNotFound,
		}
	}

	return nil
}

func (i *Implement) IsOrderItemsNotFinal(ctx context.Context, orderID int64) (bool, *exceptions.CustomError) {
	isExist, err := i.repository.IsOrderItemsNotFinal(ctx, orderID)
	if err != nil {
		return false, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check order items not final: %w", err),
		}
	}
	return isExist, nil
}

func (i *Implement) GetTableIDByOrderID(ctx context.Context, orderID int64) (int64, *exceptions.CustomError) {
	tableID, err := i.repository.GetTableIDByOrderID(ctx, orderID)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get table id by order id: %w", err),
		}
	}
	return tableID, nil
}

func (i *Implement) getOrderSequence(ctx context.Context) (string, *exceptions.CustomError) {

	currentTimeDB, err := i.repository.GetTimeNow(ctx)
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get current time: %w", err),
		}
	}

	if !currentTimeDB.Valid {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("validation failed: current time is not valid"),
		}
	}

	currentLocation, err := time.LoadLocation(i.config.TimeZone)
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: fmt.Errorf("failed to load time zone: %w", err),
		}
	}

	currentTime := currentTimeDB.Time.In(currentLocation)
	sequence, err := i.repository.GetOrderSequence(ctx, pgtype.Date{
		Time:  currentTime,
		Valid: true,
	})
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get or create order sequence: %w", err),
		}
	}

	return fmt.Sprintf("FS-%s-%04d", currentTime.In(currentLocation).Format("20060102"), sequence), nil
}
