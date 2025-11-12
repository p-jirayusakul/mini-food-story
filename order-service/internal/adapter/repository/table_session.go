package repository

import (
	"context"
	"fmt"
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

func (i *Implement) GetSessionIDByTableID(ctx context.Context, tableID int64) (result uuid.UUID, customError *exceptions.CustomError) {
	sessionID, err := i.repository.GetSessionIDByTableID(ctx, tableID)
	if err != nil {
		return uuid.UUID{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	v, err := utils.PareStringToUUID(sessionID.String())
	if err != nil {
		return uuid.UUID{}, &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: fmt.Errorf("failed to parse session id: %w", err),
		}
	}

	return v, nil
}
