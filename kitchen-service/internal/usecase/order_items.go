package usecase

import (
	"context"
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
	shareModel "food-story/shared/model"
)

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload shareModel.OrderItemsStatus) (customError *exceptions.CustomError) {
	return i.repository.UpdateOrderItemsStatus(ctx, payload)
}

func (i *Implement) UpdateOrderItemsStatusServed(ctx context.Context, payload shareModel.OrderItemsStatus) (customError *exceptions.CustomError) {
	return i.repository.UpdateOrderItemsStatusServed(ctx, payload)
}

func (i *Implement) SearchOrderItems(ctx context.Context, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, customError *exceptions.CustomError) {
	return i.repository.SearchOrderItems(ctx, payload)
}

func (i *Implement) GetOrderItems(ctx context.Context, orderID int64, search domain.SearchOrderItems) (result domain.SearchOrderItemsResult, customError *exceptions.CustomError) {
	return i.repository.GetOrderItems(ctx, orderID, search)
}

func (i *Implement) GetOrderItemsByID(ctx context.Context, orderID, orderItemsID int64) (result *shareModel.OrderItems, customError *exceptions.CustomError) {
	tableNumber, customError := i.repository.GetTableNumberOrderByID(ctx, orderID)
	if customError != nil {
		return nil, customError
	}

	return i.repository.GetOrderItemsByID(ctx, orderID, orderItemsID, tableNumber)
}
