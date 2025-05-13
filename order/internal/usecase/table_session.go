package usecase

import (
	"context"
	"food-story/pkg/exceptions"
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
