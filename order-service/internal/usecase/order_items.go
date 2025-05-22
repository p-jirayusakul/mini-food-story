package usecase

import (
	"context"
	"errors"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"
	"github.com/google/uuid"
)

func (i *Implement) CreateOrderItems(ctx context.Context, sessionID uuid.UUID, items []shareModel.OrderItems) (customError *exceptions.CustomError) {
	tableSession, customError := i.GetCurrentTableSession(sessionID)
	if customError != nil {
		return customError
	}

	orderID, convertErr := convertOrderID(*tableSession.OrderID)
	if convertErr != nil {
		return convertErr
	}

	for index := range items {
		items[index].OrderID = orderID
	}

	orderItems, customError := i.repository.CreateOrderItems(ctx, items)
	if customError != nil {
		return
	}

	customError = i.PublishOrderToQueue(orderItems)
	if customError != nil {
		return customError
	}

	return
}

func (i *Implement) GetCurrentOrderItems(ctx context.Context, sessionID uuid.UUID, pageNumber int64) (result domain.SearchCurrentOrderItemsResult, customError *exceptions.CustomError) {
	tableSession, customError := i.GetCurrentTableSession(sessionID)
	if customError != nil {
		return domain.SearchCurrentOrderItemsResult{}, customError
	}

	orderID, convertErr := convertOrderID(*tableSession.OrderID)
	if convertErr != nil {
		return domain.SearchCurrentOrderItemsResult{}, convertErr
	}

	return i.repository.GetCurrentOrderItems(ctx, orderID, pageNumber)
}

func (i *Implement) GetCurrentOrderItemsByID(ctx context.Context, sessionID uuid.UUID, orderItemsID int64) (result *domain.CurrentOrderItems, customError *exceptions.CustomError) {
	tableSession, customError := i.GetCurrentTableSession(sessionID)
	if customError != nil {
		return nil, customError
	}

	orderID, convertErr := convertOrderID(*tableSession.OrderID)
	if convertErr != nil {
		return nil, convertErr
	}

	return i.repository.GetCurrentOrderItemsByID(ctx, orderID, orderItemsID)
}

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, sessionID uuid.UUID, payload shareModel.OrderItemsStatus) (customError *exceptions.CustomError) {
	orderID, customError := i.GetOrderIDFromSession(sessionID)
	if customError != nil {
		return customError
	}

	payload.OrderID = orderID

	return i.repository.UpdateOrderItemsStatus(ctx, payload)
}

func (i *Implement) SearchOrderItemsIncomplete(ctx context.Context, orderID int64, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, customError *exceptions.CustomError) {
	return i.repository.SearchOrderItemsIncomplete(ctx, orderID, payload)
}

func convertOrderID(orderID string) (int64, *exceptions.CustomError) {
	result, err := utils.StrToInt64(orderID)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: err,
		}
	}

	if result == 0 {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: errors.New("invalid order ID"),
		}
	}

	return result, nil
}
