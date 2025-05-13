package usecase

import (
	"context"
	"errors"
	"food-story/order/internal/domain"
	"food-story/pkg/exceptions"
	"github.com/google/uuid"
)

func (i *Implement) CreateOrderItems(ctx context.Context, sessionID uuid.UUID, items []domain.OrderItems) (customError *exceptions.CustomError) {
	err := i.cache.IsCachedTableExist(sessionID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrSessionNotFound,
			}
		}

		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	return i.repository.CreateOrderItems(ctx, items)
}

func (i *Implement) GetOrderItems(ctx context.Context, orderID int64) (result []*domain.OrderItems, customError *exceptions.CustomError) {
	return i.repository.GetOrderItems(ctx, orderID)
}

func (i *Implement) GetOderItemsByID(ctx context.Context, orderID, orderItemsID int64) (result *domain.OrderItems, customError *exceptions.CustomError) {
	return i.repository.GetOderItemsByID(ctx, orderID, orderItemsID)
}

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, sessionID uuid.UUID, payload domain.OrderItemsStatus) (customError *exceptions.CustomError) {
	err := i.cache.IsCachedTableExist(sessionID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrSessionNotFound,
			}
		}

		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	return i.repository.UpdateOrderItemsStatus(ctx, payload)
}
