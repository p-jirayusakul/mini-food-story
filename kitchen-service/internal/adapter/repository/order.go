package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
)

func (i *Implement) IsOrderWithItemsExists(ctx context.Context, orderID, orderItemsID int64) (customError *exceptions.CustomError) {
	isExist, err := i.repository.IsOrderWithItemsExists(ctx, database.IsOrderWithItemsExistsParams{
		OrderID:      orderID,
		OrderItemsID: orderItemsID,
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check order item exists: %w", err),
		}
	}

	if !isExist {
		return &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: exceptions.ErrOrderItemsNotFound,
		}
	}

	return nil
}

func (i *Implement) GetTableNumberOrderByID(ctx context.Context, orderID int64) (result int32, customError *exceptions.CustomError) {
	result, err := i.repository.GetTableNumberOrderByID(ctx, orderID)
	if err != nil {

		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: errors.New("table number not found"),
			}
		}

		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get table number exists: %w", err),
		}
	}

	return result, nil
}

func (i *Implement) IsOrderExist(ctx context.Context, id int64) (customError *exceptions.CustomError) {
	isExist, err := i.repository.IsOrderExist(ctx, id)
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

func (i *Implement) IsOrderItemsNotFinal(ctx context.Context, orderID int64) (bool, *exceptions.CustomError) {
	isExist, err := i.repository.IsOrderItemsNotFinal(ctx, orderID)
	if err != nil {
		return false, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check order items not final: %w", err),
		}
	}
	return isExist, nil
}

func (i *Implement) GetTableIDByOrderID(ctx context.Context, orderID int64) (int64, *exceptions.CustomError) {
	tableID, err := i.repository.GetTableIDByOrderID(ctx, orderID)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get table id by order id: %w", err),
		}
	}
	return tableID, nil
}
