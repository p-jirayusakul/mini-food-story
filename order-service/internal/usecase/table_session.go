package usecase

import (
	"errors"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"github.com/google/uuid"
)

func (i *Implement) GetOrderIDFromSession(sessionID uuid.UUID) (result int64, customError *exceptions.CustomError) {
	session, tableCacheErr := i.cache.GetCachedTable(sessionID)
	if tableCacheErr != nil {
		return 0, tableCacheErr
	}

	if session.OrderID == nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: exceptions.ErrOrderNotFound,
		}
	}

	orderID, err := utils.StrToInt64(*session.OrderID)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: err,
		}
	}

	return orderID, nil
}

func (i *Implement) GetCurrentTableSession(sessionID uuid.UUID) (result domain.CurrentTableSession, customError *exceptions.CustomError) {
	session, tableCacheErr := i.cache.GetCachedTable(sessionID)
	if tableCacheErr != nil {
		return domain.CurrentTableSession{}, tableCacheErr
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
			Errors: exceptions.ErrOrderNotFound,
		}
	}

	return *session, nil
}
