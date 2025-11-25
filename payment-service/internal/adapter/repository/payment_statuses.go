package repository

import (
	"context"
	"food-story/pkg/exceptions"
)

func (i *Implement) GetPaymentStatusPending(ctx context.Context) (result int64, err error) {
	result, err = i.repository.GetPaymentStatusPending(ctx)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch payment status", err)
	}

	return result, nil
}
