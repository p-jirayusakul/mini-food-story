package usecase

import (
	"context"
)

func (i *PaymentImplement) UpdateTablesStatusWaitingForPayment(ctx context.Context, tableID int64) (err error) {
	err = i.repository.UpdateTablesStatusWaitingForPayment(ctx, tableID)
	if err != nil {
		return err
	}
	return nil
}
