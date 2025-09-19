package usecase

import (
	"context"
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
	shareModel "food-story/shared/model"
)

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload shareModel.OrderItemsStatus) (customError *exceptions.CustomError) {

	customError = i.repository.UpdateOrderItemsStatus(ctx, payload)
	if customError != nil {
		return customError
	}

	return i.updateTablesStatusFoodServed(ctx, payload.OrderID)
}

func (i *Implement) UpdateOrderItemsStatusServed(ctx context.Context, payload shareModel.OrderItemsStatus) (customError *exceptions.CustomError) {
	customError = i.repository.UpdateOrderItemsStatusServed(ctx, payload)
	if customError != nil {
		return customError
	}

	return i.updateTablesStatusFoodServed(ctx, payload.OrderID)
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

func (i *Implement) updateTablesStatusFoodServed(ctx context.Context, orderID int64) (customError *exceptions.CustomError) {
	isOrderItemsNotFinal, customError := i.repository.IsOrderItemsNotFinal(ctx, orderID)
	if customError != nil {
		return customError
	}

	if !isOrderItemsNotFinal {
		tableID, tableIDErr := i.repository.GetTableIDByOrderID(ctx, orderID)
		if tableIDErr != nil {
			return tableIDErr
		}

		statusFoodServedErr := i.repository.UpdateTablesStatusFoodServed(ctx, tableID)
		if statusFoodServedErr != nil {
			return statusFoodServedErr
		}
	}

	return
}
