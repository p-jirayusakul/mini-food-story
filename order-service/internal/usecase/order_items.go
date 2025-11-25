package usecase

import (
	"context"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"

	"github.com/google/uuid"
)

func (i *Implement) CreateOrderItems(ctx context.Context, sessionID uuid.UUID, items []shareModel.OrderItems) (err error) {
	tableSession, err := i.GetCurrentTableSession(sessionID)
	if err != nil {
		return err
	}

	orderID, err := convertOrderID(*tableSession.OrderID)
	if err != nil {
		return err
	}

	for index := range items {
		items[index].OrderID = orderID
	}

	orderItems, err := i.repository.CreateOrderItems(ctx, items)
	if err != nil {
		return
	}

	err = i.PublishOrderToQueue(orderItems)
	if err != nil {
		return err
	}

	return
}

func (i *Implement) GetCurrentOrderItems(ctx context.Context, sessionID uuid.UUID, pageNumber, pageSize int64) (result domain.SearchCurrentOrderItemsResult, err error) {
	tableSession, err := i.GetCurrentTableSession(sessionID)
	if err != nil {
		return domain.SearchCurrentOrderItemsResult{}, err
	}

	orderID, err := convertOrderID(*tableSession.OrderID)
	if err != nil {
		return domain.SearchCurrentOrderItemsResult{}, err
	}

	return i.repository.GetCurrentOrderItems(ctx, orderID, pageNumber, pageSize)
}

func (i *Implement) GetCurrentOrderItemsByID(ctx context.Context, sessionID uuid.UUID, orderItemsID int64) (result *domain.CurrentOrderItems, err error) {
	tableSession, err := i.GetCurrentTableSession(sessionID)
	if err != nil {
		return nil, err
	}

	orderID, err := convertOrderID(*tableSession.OrderID)
	if err != nil {
		return nil, err
	}

	return i.repository.GetCurrentOrderItemsByID(ctx, orderID, orderItemsID)
}

func (i *Implement) GetOrderItems(ctx context.Context, orderID, pageNumber, pageSize int64) (result domain.SearchCurrentOrderItemsResult, err error) {
	return i.repository.GetCurrentOrderItems(ctx, orderID, pageNumber, pageSize)
}

func (i *Implement) UpdateOrderItemsStatusByID(ctx context.Context, payload shareModel.OrderItemsStatus) (err error) {
	err = i.repository.UpdateOrderItemsStatus(ctx, payload)
	if err != nil {
		return err
	}

	isOrderItemsNotFinal, err := i.repository.IsOrderItemsNotFinal(ctx, payload.OrderID)
	if err != nil {
		return err
	}

	if !isOrderItemsNotFinal {
		tableID, tableIDErr := i.repository.GetTableIDByOrderID(ctx, payload.OrderID)
		if tableIDErr != nil {
			return tableIDErr
		}

		updateErr := i.repository.UpdateTablesStatusFoodServed(ctx, tableID)
		if updateErr != nil {
			return updateErr
		}
	}

	return
}

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, sessionID uuid.UUID, payload shareModel.OrderItemsStatus) (err error) {
	orderID, err := i.GetOrderIDFromSession(sessionID)
	if err != nil {
		return err
	}

	payload.OrderID = orderID

	return i.UpdateOrderItemsStatusByID(ctx, payload)
}

func (i *Implement) SearchOrderItemsIncomplete(ctx context.Context, orderID int64, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, err error) {
	return i.repository.SearchOrderItemsIncomplete(ctx, orderID, payload)
}

func convertOrderID(orderID string) (int64, error) {
	result, err := utils.StrToInt64(orderID)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeSystem, "failed to convert order ID to int64", err)
	}

	if result == 0 {
		return 0, exceptions.Error(exceptions.CodeBusiness, "order ID cannot be zero")
	}

	return result, nil
}
