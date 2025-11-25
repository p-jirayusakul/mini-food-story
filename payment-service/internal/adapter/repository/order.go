package repository

import (
	"context"
	"food-story/pkg/exceptions"

	"github.com/google/uuid"
)

func (i *Implement) IsOrderItemsNotFinal(ctx context.Context, orderID int64) (err error) {
	isOrderItemsNotFinal, err := i.repository.IsOrderItemsNotFinal(ctx, orderID)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to fetch order item not final", err)
	}

	if isOrderItemsNotFinal {
		return exceptions.Error(exceptions.CodeBusiness, "all order items is not 'served' or 'cancel', cannot be paid")
	}

	return nil
}

func (i *Implement) IsOrderExist(ctx context.Context, orderID int64) (err error) {
	isExist, err := i.repository.IsOrderExist(ctx, orderID)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to check order exists", err)
	}

	if !isExist {
		return exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderNotFound.Error())
	}

	return nil
}
func (i *Implement) GetTableIDByOrderID(ctx context.Context, orderID int64) (result int64, err error) {
	tableID, err := i.repository.GetTableIDByOrderID(ctx, orderID)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to get table by order id", err)
	}

	return tableID, nil
}

func (i *Implement) GetSessionIDByOrderID(ctx context.Context, orderID int64) (result uuid.UUID, err error) {
	sessionIDData, err := i.repository.GetSessionIDByOrderID(ctx, orderID)
	if err != nil {
		return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, "failed to get session by order id", err)
	}

	sessionIDString := sessionIDData.String()
	sessionID, err := uuid.Parse(sessionIDString)
	if err != nil {
		return uuid.Nil, exceptions.Errorf(exceptions.CodeSystem, "failed to get session by order id", err)
	}

	return sessionID, nil
}
