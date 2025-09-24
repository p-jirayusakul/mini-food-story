package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
)

func (i *Implement) IsOrderItemsNotFinal(ctx context.Context, orderID int64) (customError *exceptions.CustomError) {
	isOrderItemsNotFinal, err := i.repository.IsOrderItemsNotFinal(ctx, orderID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	if isOrderItemsNotFinal {
		return &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: errors.New("all order items is not 'served' or 'cancel', cannot be paid"),
		}
	}

	return nil
}

func (i *Implement) IsOrderExist(ctx context.Context, orderID int64) (customError *exceptions.CustomError) {
	isExist, err := i.repository.IsOrderExist(ctx, orderID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check order exists: %w", err),
		}
	}

	if !isExist {
		return &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: exceptions.ErrOrderNotFound,
		}
	}

	return nil
}

func (i *Implement) GetTableIDByOrderID(ctx context.Context, orderID int64) (result int64, customError *exceptions.CustomError) {
	tableID, err := i.repository.GetTableIDByOrderID(ctx, orderID)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get table by order id: %w", err),
		}
	}

	return tableID, nil
}
