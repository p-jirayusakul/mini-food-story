package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"

	"github.com/google/uuid"
)

func (i *Implement) GetOrderIDBySessionID(ctx context.Context, sessionID uuid.UUID) (orderID int64, customError *exceptions.CustomError) {
	data, err := i.repository.GetOrderIDBySessionID(ctx, utils.UUIDToPgUUID(sessionID))
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: errors.New("order id not found"),
			}
		}
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch order id by session id: %w", err),
		}
	}
	return data, nil
}
