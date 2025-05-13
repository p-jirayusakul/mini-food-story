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

	return i.repository.CreateOrder(ctx, order)
}

func (i *Implement) GetOrderByID(ctx context.Context, id int64) (result *domain.Order, customError *exceptions.CustomError) {
	return i.repository.GetOrderByID(ctx, id)
}

func (i *Implement) UpdateOrderStatus(ctx context.Context, sessionID uuid.UUID, payload domain.OrderStatus) (customError *exceptions.CustomError) {
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

	customError = i.repository.UpdateOrderStatus(ctx, payload)
	if customError != nil {
		return
	}

	return i.UpdateOrderStatusClosed(ctx, sessionID, payload.StatusCode)
}
