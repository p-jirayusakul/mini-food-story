package usecase

import (
	"context"
	"fmt"
	"food-story/order/internal/domain"
	"food-story/pkg/exceptions"
	"github.com/google/uuid"
)

func (i *Implement) CreateOrderItems(ctx context.Context, sessionID uuid.UUID, items []domain.OrderItems) (customError *exceptions.CustomError) {
	orderID, customError := i.GetOrderIDFromSession(sessionID)
	if customError != nil {
		return
	}

	for index, _ := range items {
		items[index].OrderID = orderID
	}

	return i.repository.CreateOrderItems(ctx, items)
}

func (i *Implement) GetOrderItems(ctx context.Context, sessionID uuid.UUID) (result []*domain.OrderItems, customError *exceptions.CustomError) {
	orderID, customError := i.GetOrderIDFromSession(sessionID)
	if customError != nil {
		return
	}

	return i.repository.GetOrderItems(ctx, orderID)
}

func (i *Implement) GetOderItemsByID(ctx context.Context, sessionID uuid.UUID, orderItemsID int64) (result *domain.OrderItems, customError *exceptions.CustomError) {
	orderID, customError := i.GetOrderIDFromSession(sessionID)
	if customError != nil {
		return
	}

	return i.repository.GetOderItemsByID(ctx, orderID, orderItemsID)
}

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, sessionID uuid.UUID, payload domain.OrderItemsStatus) (customError *exceptions.CustomError) {
	orderID, customError := i.GetOrderIDFromSession(sessionID)
	if customError != nil {
		return
	}

	payload.OrderID = orderID

	fmt.Println(payload.OrderID)

	return i.repository.UpdateOrderItemsStatus(ctx, payload)
}
