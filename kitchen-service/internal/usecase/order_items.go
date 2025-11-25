package usecase

import (
	"context"
	"food-story/kitchen-service/internal/domain"
	shareModel "food-story/shared/model"
)

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload shareModel.OrderItemsStatus) (err error) {

	err = i.repository.UpdateOrderItemsStatus(ctx, payload)
	if err != nil {
		return err
	}

	return i.updateTablesStatusFoodServed(ctx, payload.OrderID)
}

func (i *Implement) UpdateOrderItemsStatusServed(ctx context.Context, payload shareModel.OrderItemsStatus) (err error) {
	err = i.repository.UpdateOrderItemsStatusServed(ctx, payload)
	if err != nil {
		return err
	}

	return i.updateTablesStatusFoodServed(ctx, payload.OrderID)
}

func (i *Implement) SearchOrderItems(ctx context.Context, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, err error) {
	return i.repository.SearchOrderItems(ctx, payload)
}

func (i *Implement) GetOrderItems(ctx context.Context, orderID int64, search domain.SearchOrderItems) (result domain.SearchOrderItemsResult, err error) {
	return i.repository.GetOrderItems(ctx, orderID, search)
}

func (i *Implement) GetOrderItemsByID(ctx context.Context, orderID, orderItemsID int64) (result *shareModel.OrderItems, err error) {
	tableNumber, err := i.repository.GetTableNumberOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return i.repository.GetOrderItemsByID(ctx, orderID, orderItemsID, tableNumber)
}

func (i *Implement) updateTablesStatusFoodServed(ctx context.Context, orderID int64) (err error) {
	isOrderItemsNotFinal, err := i.repository.IsOrderItemsNotFinal(ctx, orderID)
	if err != nil {
		return err
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
