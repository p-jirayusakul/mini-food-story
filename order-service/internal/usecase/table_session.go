package usecase

import (
	"context"
	"errors"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"

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

func (i *Implement) GetCurrentTableSession(sessionID uuid.UUID) (result shareModel.CurrentTableSession, customError *exceptions.CustomError) {
	session, tableCacheErr := i.cache.GetCachedTable(sessionID)
	if tableCacheErr != nil {
		return shareModel.CurrentTableSession{}, tableCacheErr
	}

	if session == nil {
		return shareModel.CurrentTableSession{}, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: errors.New("session not found"),
		}
	}

	if session.OrderID == nil {
		return shareModel.CurrentTableSession{}, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: exceptions.ErrOrderNotFound,
		}
	}

	return *session, nil
}

func (i *Implement) GetSessionIDByTableID(ctx context.Context, tableID int64) (result uuid.UUID, customError *exceptions.CustomError) {
	return i.repository.GetSessionIDByTableID(ctx, tableID)
}

func (i *Implement) IsSessionValid(sessionID uuid.UUID) *exceptions.CustomError {
	return i.cache.IsCachedTableExist(sessionID)
}
