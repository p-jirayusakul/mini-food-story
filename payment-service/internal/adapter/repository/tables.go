package repository

import (
	"context"
	"food-story/pkg/exceptions"
)

func (i *Implement) UpdateTablesStatusWaitingForPayment(ctx context.Context, tableID int64) (err error) {
	err = i.repository.UpdateTablesStatusWaitingForPayment(ctx, tableID)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to update table status waiting for payment", err)
	}

	return nil
}

func (i *Implement) UpdateTablesStatusCleaning(ctx context.Context, tableID int64) (err error) {
	err = i.repository.UpdateTablesStatusCleaning(ctx, tableID)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to update table status cleaning", err)
	}

	return nil
}
