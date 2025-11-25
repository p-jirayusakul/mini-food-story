package repository

import (
	"context"
	"food-story/pkg/exceptions"
)

func (i *Implement) IsOrderStatusExist(ctx context.Context, statusCode string) (err error) {
	isStatusExist, err := i.repository.IsOrderStatusExist(ctx, statusCode)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to check order status exists", err)
	}

	if !isStatusExist {
		return exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderStatusNotFound.Error())
	}

	return nil
}

func (i *Implement) IsOrderStatusFinal(ctx context.Context, statusCode string) (result bool, err error) {
	isFinalStatus, err := i.repository.IsOrderStatusFinal(ctx, statusCode)
	if err != nil {
		return false, exceptions.Errorf(exceptions.CodeRepository, "failed to check order status exists", err)
	}

	return isFinalStatus, nil
}

func (i *Implement) GetOrderStatusPreparing(ctx context.Context) (result int64, err error) {
	id, err := i.repository.GetOrderStatusPreparing(ctx)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to get order status preparing", err)
	}

	return id, nil
}
