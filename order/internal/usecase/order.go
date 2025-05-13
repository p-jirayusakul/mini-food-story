package usecase

import (
	"context"
	"errors"
	"food-story/order/internal/domain"
	"food-story/pkg/exceptions"
	"github.com/google/uuid"
)

func (i *Implement) CreateOrder(ctx context.Context, sessionID uuid.UUID) (result int64, customError *exceptions.CustomError) {
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

	order := domain.Order{
		SessionID: sessionID,
		TableID:   sessionDetail.TableID,
	}

	orderID, customError := i.repository.CreateOrder(ctx, order)
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
