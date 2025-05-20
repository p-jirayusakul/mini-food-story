package repository

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"github.com/google/uuid"
)

func (i *Implement) UpdateStatusCloseTableSession(ctx context.Context, sessionID uuid.UUID) (customError *exceptions.CustomError) {
	err := i.repository.UpdateStatusCloseTableSession(ctx, utils.UUIDToPgUUID(sessionID))
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	return nil
}
