package repository

import (
	"context"
	"fmt"
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
)

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload domain.OrderItemsStatus) (customError *exceptions.CustomError) {
	customError = i.IsOrderExist(ctx, payload.OrderID)
	if customError != nil {
		return
	}

	customError = i.IsOrderStatus(ctx, payload.StatusCode)
	if customError != nil {
		return
	}

	err := i.repository.UpdateOrderItemsStatus(ctx, database.UpdateOrderItemsStatusParams{
		StatusCode: payload.StatusCode,
		ID:         payload.ID,
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update order items status: %w", err),
		}
	}

	return
}
