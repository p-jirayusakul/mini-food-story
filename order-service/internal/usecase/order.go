package usecase

import (
	"context"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"

	"github.com/google/uuid"
)

func (i *Implement) CreateOrder(ctx context.Context, sessionID uuid.UUID, orderItems []shareModel.OrderItems) (result int64, err error) {

	if len(orderItems) == 0 {
		return 0, exceptions.Error(exceptions.CodeBusiness, "order items cannot be empty")
	}

	if sessionID == uuid.Nil {
		return 0, exceptions.Error(exceptions.CodeBusiness, "invalid table ID")
	}

	sessionDetail, tableCacheErr := i.cache.GetCachedTable(sessionID)
	if tableCacheErr != nil {
		return 0, tableCacheErr
	}

	// if already have order ID, then create order items
	if sessionDetail.OrderID != nil {
		orderID, sysErr := utils.StrToInt64(*sessionDetail.OrderID)
		if sysErr != nil {
			return 0, exceptions.Error(exceptions.CodeSystem, sysErr.Error())
		}

		if len(orderItems) > 0 {
			createErr := i.CreateOrderItems(ctx, sessionID, orderItems)
			if createErr != nil {
				return 0, createErr
			}
		}

		return orderID, nil
	}

	// if not have order ID, then create order
	payloadCreateOrder := domain.CreateOrder{
		SessionID:  sessionID,
		TableID:    sessionDetail.TableID,
		OrderItems: orderItems,
	}

	orderID, err := i.repository.CreateOrder(ctx, payloadCreateOrder)
	if err != nil {
		return 0, err
	}

	err = i.cache.UpdateOrderID(sessionID, orderID)
	if err != nil {
		return 0, err
	}

	newOrderItems, err := i.repository.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return 0, err
	}

	err = i.PublishOrderToQueue(newOrderItems)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (i *Implement) GetOrderByID(ctx context.Context, sessionID uuid.UUID) (result *domain.Order, err error) {
	orderID, err := i.GetOrderIDFromSession(sessionID)
	if err != nil {
		return nil, err
	}

	return i.repository.GetOrderByID(ctx, orderID)
}

func (i *Implement) GetSessionIDByOrderID(ctx context.Context, orderID int64) (result uuid.UUID, err error) {
	return i.repository.GetSessionIDByOrderID(ctx, orderID)
}
