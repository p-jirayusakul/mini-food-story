package repository

import (
	"context"
	"food-story/pkg/exceptions"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (i *Implement) UpdateStatusCloseTableSession(ctx context.Context, sessionID uuid.UUID) (customError *exceptions.CustomError) {
	var sessionByte [16]byte = sessionID
	err := i.repository.UpdateStatusCloseTableSession(ctx, pgtype.UUID{Bytes: sessionByte, Valid: true})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	return nil
}
