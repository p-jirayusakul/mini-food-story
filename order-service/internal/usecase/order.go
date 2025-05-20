package usecase

import (
	"context"
	"fmt"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"
	"github.com/google/uuid"
)

func (i *Implement) CreateOrder(ctx context.Context, sessionID uuid.UUID, orderItems []shareModel.OrderItems) (result int64, customError *exceptions.CustomError) {

	if len(orderItems) == 0 {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: fmt.Errorf("order items cannot be empty"),
		}
	}

	if sessionID == uuid.Nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: fmt.Errorf("invalid table ID"),
		}
	}

	sessionDetail, tableCacheErr := i.cache.GetCachedTable(sessionID)
	if tableCacheErr != nil {
		return 0, tableCacheErr
	}

	if sessionDetail.OrderID != nil {
		orderID, err := utils.StrToInt64(*sessionDetail.OrderID)
		if err != nil {
			return 0, &exceptions.CustomError{
				Status: exceptions.ERRSYSTEM,
				Errors: err,
			}
		}

		if len(orderItems) > 0 {
			createErr := i.CreateOrderItems(ctx, sessionID, orderItems)
			if createErr != nil {
				return 0, createErr
			}
		}

		return orderID, nil
	}

	payloadCreateOrder := domain.CreateOrder{
		SessionID:  sessionID,
		TableID:    sessionDetail.TableID,
		OrderItems: orderItems,
	}

	orderID, customError := i.repository.CreateOrder(ctx, payloadCreateOrder)
	if customError != nil {
		return 0, customError
	}

	updateCacheErr := i.cache.UpdateOrderID(sessionID, orderID)
	if updateCacheErr != nil {
		return 0, updateCacheErr
	}

	newOrderItems, getOrderItemsErr := i.repository.GetOrderItems(ctx, orderID)
	if getOrderItemsErr != nil {
		return 0, getOrderItemsErr
	}

	queueErr := i.PublishOrderToQueue(newOrderItems)
	if queueErr != nil {
		return 0, queueErr
	}

	return orderID, nil
}

func (i *Implement) GetOrderByID(ctx context.Context, sessionID uuid.UUID) (result *domain.Order, customError *exceptions.CustomError) {
	orderID, customError := i.GetOrderIDFromSession(sessionID)
	if customError != nil {
		return nil, customError
	}

	return i.repository.GetOrderByID(ctx, orderID)
}
