package usecase

import (
	"context"
	"errors"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"github.com/google/uuid"
)

func (i *Implement) UpdateOrderStatusClosed(ctx context.Context, sessionID uuid.UUID, statusCode string) (customError *exceptions.CustomError) {
	isFinalStatus, customError := i.repository.IsOrderStatusFinal(ctx, statusCode)
	if customError != nil {
		return
	}

	if isFinalStatus {
		err := i.cache.DeleteCachedTable(sessionID)
		if err != nil {
			return &exceptions.CustomError{
				Status: exceptions.ERRREPOSITORY,
				Errors: err,
			}
		}

		customError = i.repository.UpdateStatusCloseTableSession(ctx, sessionID)
		if customError != nil {
			return customError
		}

		return nil
	}

	return nil
}

func (i *Implement) GetOrderIDFromSession(sessionID uuid.UUID) (result int64, customError *exceptions.CustomError) {
	session, err := i.cache.GetCachedTable(sessionID)
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

	if session.OrderID == nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: errors.New("order not found"),
		}
	}

	orderID, err := utils.StrToInt64(*session.OrderID)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: err,
		}
	}

	return orderID, nil
}

func (i *Implement) GetCurrentTableSession(sessionID uuid.UUID) (result domain.CurrentTableSession, customError *exceptions.CustomError) {
	session, err := i.cache.GetCachedTable(sessionID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFound) {
			return domain.CurrentTableSession{}, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrSessionNotFound,
			}
		}

		return domain.CurrentTableSession{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	if session == nil {
		return domain.CurrentTableSession{}, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: errors.New("session not found"),
		}
	}

	if session.OrderID == nil {
		return domain.CurrentTableSession{}, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: errors.New("order not found"),
		}
	}

	return *session, nil
}
