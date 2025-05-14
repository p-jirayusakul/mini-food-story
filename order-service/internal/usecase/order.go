package usecase

import (
	"context"
	"errors"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"github.com/google/uuid"
)

func (i *Implement) CreateOrder(ctx context.Context, sessionID uuid.UUID, items []domain.OrderItems) (result int64, customError *exceptions.CustomError) {
	sessionDetail, err := i.cache.GetCachedTable(sessionID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return 0, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrSessionNotFound,
			}
		}

		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	if sessionDetail.OrderID != nil {
		orderID, err := utils.StrToInt64(*sessionDetail.OrderID)
		if err != nil {
			return 0, &exceptions.CustomError{
				Status: exceptions.ERRUNKNOWN,
				Errors: err,
			}
		}

		if len(items) > 0 {
			customError := i.CreateOrderItems(ctx, sessionID, items)
			if customError != nil {
				return 0, customError
			}
		}

		return orderID, nil
	}

	payloadCreateOrder := domain.CreateOrder{
		SessionID:  sessionID,
		TableID:    sessionDetail.TableID,
		OrderItems: items,
	}

	orderID, customError := i.repository.CreateOrder(ctx, payloadCreateOrder)
	if customError != nil {
		return 0, customError
	}

	err = i.cache.UpdateOrderID(sessionID, orderID)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	// public message to kafka
	orderItems, customError := i.GetOrderItems(ctx, sessionID)
	if customError != nil {
		return 0, customError
	}
	if len(orderItems) > 0 {
		for _, item := range orderItems {
			err := i.queue.PublishOrder(*item)
			if err != nil {
				return 0, &exceptions.CustomError{
					Status: exceptions.ERRREPOSITORY,
					Errors: err,
				}
			}
		}
	}

	return orderID, nil
}

func (i *Implement) GetOrderByID(ctx context.Context, sessionID uuid.UUID) (result *domain.Order, customError *exceptions.CustomError) {
	orderID, customError := i.GetOrderIDFromSession(sessionID)
	if customError != nil {
		return
	}

	return i.repository.GetOrderByID(ctx, orderID)
}

func (i *Implement) UpdateOrderStatus(ctx context.Context, sessionID uuid.UUID, payload domain.OrderStatus) (customError *exceptions.CustomError) {
	orderID, customError := i.GetOrderIDFromSession(sessionID)
	if customError != nil {
		return
	}

	payload.ID = orderID

	customError = i.repository.UpdateOrderStatus(ctx, payload)
	if customError != nil {
		return
	}

	return i.UpdateOrderStatusClosed(ctx, sessionID, payload.StatusCode)
}
