package usecase

import (
	"context"
	"food-story/pkg/exceptions"
)

func (i *PaymentImplement) UpdateTablesStatusWaitingForPayment(ctx context.Context, tableID int64) (customError *exceptions.CustomError) {
	customError = i.repository.UpdateTablesStatusWaitingForPayment(ctx, tableID)
	if customError != nil {
		return customError
	}
	return nil
}
