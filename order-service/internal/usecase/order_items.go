package usecase

import (
	"context"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"github.com/google/uuid"
)

func (i *Implement) CreateOrderItems(ctx context.Context, sessionID uuid.UUID, items []domain.OrderItems) (customError *exceptions.CustomError) {
	tableSession, customError := i.GetCurrentTableSession(sessionID)
	if customError != nil {
		return
	}

	orderID, err := utils.StrToInt64(*tableSession.OrderID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: err,
		}
	}

	for index := range items {
		items[index].OrderID = orderID
	}

	orderItems, customError := i.repository.CreateOrderItems(ctx, items, tableSession.TableNumber)
	if customError != nil {
		return
	}

	customError = i.PublishOrderToQueue(orderItems)
	if customError != nil {
		return customError
	}

	return
}

func (i *Implement) GetOrderItems(ctx context.Context, sessionID uuid.UUID) (result []*domain.OrderItems, customError *exceptions.CustomError) {
	tableSession, customError := i.GetCurrentTableSession(sessionID)
	if customError != nil {
		return
	}

	orderID, err := utils.StrToInt64(*tableSession.OrderID)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: err,
		}
	}

	return i.repository.GetOrderItems(ctx, orderID, tableSession.TableNumber)
}

func (i *Implement) GetOderItemsByID(ctx context.Context, sessionID uuid.UUID, orderItemsID int64) (result *domain.OrderItems, customError *exceptions.CustomError) {
	tableSession, customError := i.GetCurrentTableSession(sessionID)
	if customError != nil {
		return
	}

	orderID, err := utils.StrToInt64(*tableSession.OrderID)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: err,
		}
	}

	return i.repository.GetOderItemsByID(ctx, orderID, orderItemsID, tableSession.TableNumber)
}

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, sessionID uuid.UUID, payload domain.OrderItemsStatus) (customError *exceptions.CustomError) {
	orderID, customError := i.GetOrderIDFromSession(sessionID)
	if customError != nil {
		return
	}

	payload.OrderID = orderID

	return i.repository.UpdateOrderItemsStatus(ctx, payload)
}

func (i *Implement) SearchOrderItemsIncomplete(ctx context.Context, orderID int64, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, customError *exceptions.CustomError) {
	return i.repository.SearchOrderItemsIncomplete(ctx, orderID, payload)
}
