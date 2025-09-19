package repository

import (
	"context"
	"fmt"
	"food-story/pkg/exceptions"
)

func (i *Implement) UpdateTablesStatusFoodServed(ctx context.Context, tableID int64) (customError *exceptions.CustomError) {
	err := i.repository.UpdateTablesStatusFoodServed(ctx, tableID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update table status: %w", err),
		}
	}

	return nil
}

func (i *Implement) UpdateTablesStatusWaitingToBeServed(ctx context.Context, tableID int64) (customError *exceptions.CustomError) {
	err := i.repository.UpdateTablesStatusWaitingToBeServed(ctx, tableID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update table status: %w", err),
		}
	}

	return nil
}
