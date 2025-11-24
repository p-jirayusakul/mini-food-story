package repository

import (
	"context"
	"errors"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"

	"github.com/google/uuid"
)

func (i *Implement) GetOrderIDBySessionID(ctx context.Context, sessionID uuid.UUID) (orderID int64, err error) {
	data, err := i.repository.GetOrderIDBySessionID(ctx, utils.UUIDToPgUUID(sessionID))
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderNotFound.Error())
		}
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch order id by session id", err)
	}
	return data, nil
}
