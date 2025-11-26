package repository

import (
	"context"
	"food-story/pkg/exceptions"
)

func (i *Implement) GetOrderStatusPreparing(ctx context.Context) (result int64, err error) {
	id, err := i.repository.GetOrderStatusPreparing(ctx)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to get order status preparing", err)
	}

	return id, nil
}
