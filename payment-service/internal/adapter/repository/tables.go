package repository

import (
	"context"
	"fmt"
	"food-story/pkg/exceptions"
)

func (i *Implement) UpdateTablesStatusWaitingForPayment(ctx context.Context, tableID int64) (customError *exceptions.CustomError) {
	err := i.repository.UpdateTablesStatusWaitingForPayment(ctx, tableID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update table status waiting for payment: %w", err),
		}
	}

	return nil
}

func (i *Implement) UpdateTablesStatusCleaning(ctx context.Context, tableID int64) (customError *exceptions.CustomError) {
	err := i.repository.UpdateTablesStatusCleaning(ctx, tableID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update table status cleaning: %w", err),
		}
	}

	return nil
}
