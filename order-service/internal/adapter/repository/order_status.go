package repository

import (
	"context"
	"fmt"
	"food-story/pkg/exceptions"
)

func (i *Implement) IsOrderStatusExist(ctx context.Context, statusCode string) (customError *exceptions.CustomError) {
	isStatusExist, err := i.repository.IsOrderStatusExist(ctx, statusCode)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check order status exists: %w", err),
		}
	}

	if !isStatusExist {
		return &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: fmt.Errorf("order status not found"),
		}
	}

	return nil
}

func (i *Implement) GetOrderStatusPreparing(ctx context.Context) (result int64, customError *exceptions.CustomError) {
	id, err := i.repository.GetOrderStatusPreparing(ctx)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get order status preparing: %w", err),
		}
	}

	return id, nil
}
