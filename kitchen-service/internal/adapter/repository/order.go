package repository

import (
	"context"
	"errors"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
)

func (i *Implement) IsOrderWithItemsExists(ctx context.Context, orderID, orderItemsID int64) (err error) {
	isExist, err := i.repository.IsOrderWithItemsExists(ctx, database.IsOrderWithItemsExistsParams{
		OrderID:      orderID,
		OrderItemsID: orderItemsID,
	})
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to check order item exists", err)
	}

	if !isExist {
		return exceptions.ErrorIDNotFound(exceptions.CodeOrderNotFound, orderID)
	}

	return nil
}
func (i *Implement) GetTableNumberOrderByID(ctx context.Context, orderID int64) (result int32, err error) {
	result, err = i.repository.GetTableNumberOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, exceptions.ErrorIDNotFound(exceptions.CodeTableNotFound, 0)
		}
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to get table number exists", err)
	}

	return result, nil
}

func (i *Implement) IsOrderExist(ctx context.Context, id int64) (err error) {
	isExist, err := i.repository.IsOrderExist(ctx, id)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to check order exists", err)
	}

	if !isExist {
		return exceptions.ErrorIDNotFound(exceptions.CodeOrderNotFound, id)
	}

	return nil
}

func (i *Implement) IsOrderItemsNotFinal(ctx context.Context, orderID int64) (bool, error) {
	isExist, err := i.repository.IsOrderItemsNotFinal(ctx, orderID)
	if err != nil {
		return false, exceptions.Errorf(exceptions.CodeRepository, "failed to check order items not final", err)
	}
	return isExist, nil
}

func (i *Implement) GetTableIDByOrderID(ctx context.Context, orderID int64) (int64, error) {
	tableID, err := i.repository.GetTableIDByOrderID(ctx, orderID)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to get table id by order id", err)
	}
	return tableID, nil
}
