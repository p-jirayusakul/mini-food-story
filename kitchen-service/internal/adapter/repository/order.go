package repository

import (
	"context"
	"fmt"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
)

func (i *Implement) IsOrderWithItemsExists(ctx context.Context, orderID, orderItemsID int64) (customError *exceptions.CustomError) {
	isExist, err := i.repository.IsOrderWithItemsExists(ctx, database.IsOrderWithItemsExistsParams{
		OrderID:      orderID,
		OrderItemsID: orderItemsID,
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check order item exists: %w", err),
		}
	}

	if !isExist {
		return &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: fmt.Errorf("order item not found"),
		}
	}

	return nil
}
