package repository

import (
	"context"
	"fmt"
	"food-story/pkg/exceptions"
)

func (i *Implement) GetPaymentStatusPending(ctx context.Context) (result int64, customError *exceptions.CustomError) {
	result, err := i.repository.GetPaymentStatusPending(ctx)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch payment status: %w", err),
		}
	}

	return result, nil
}
