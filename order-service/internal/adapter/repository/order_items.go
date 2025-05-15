package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
)

func (i *Implement) CreateOrderItems(ctx context.Context, items []domain.OrderItems) (customError *exceptions.CustomError) {

	orderItems, customError := i.BuildPayloadOrderItems(ctx, items)
	if customError != nil {
		return
	}

	if len(orderItems) > 0 {
		_, err := i.repository.CreateOrderItems(ctx, orderItems)
		if err != nil {
			return &exceptions.CustomError{
				Status: exceptions.ERRREPOSITORY,
				Errors: fmt.Errorf("failed to create order items: %w", err),
			}
		}
	}

	return nil
}

func (i *Implement) BuildPayloadOrderItems(ctx context.Context, items []domain.OrderItems) ([]database.CreateOrderItemsParams, *exceptions.CustomError) {
	var orderItems []database.CreateOrderItemsParams
	for _, item := range items {
		product, err := i.repository.GetProductByID(ctx, item.ProductID)
		if err != nil {
			if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
				return []database.CreateOrderItemsParams{}, &exceptions.CustomError{
					Status: exceptions.ERRNOTFOUND,
					Errors: fmt.Errorf("product %d not found", item.ProductID),
				}
			}
			return []database.CreateOrderItemsParams{}, &exceptions.CustomError{
				Status: exceptions.ERRREPOSITORY,
				Errors: fmt.Errorf("failed to get product: %w", err),
			}
		}

		statusPreparingID, customError := i.GetOrderStatusPreparing(ctx)
		if customError != nil {
			return []database.CreateOrderItemsParams{}, customError
		}

		orderItems = append(orderItems, database.CreateOrderItemsParams{
			ID:            i.snowflakeID.Generate(),
			OrderID:       item.OrderID,
			ProductID:     product.ID,
			StatusID:      statusPreparingID,
			ProductName:   product.Name,
			ProductNameEn: product.NameEn,
			Price:         product.Price,
			Quantity:      item.Quantity,
			Note:          utils.StringPtrToPgText(item.Note),
		})
	}

	return orderItems, nil
}

func (i *Implement) GetOrderItems(ctx context.Context, orderID int64, tableNumber int32) (result []*domain.OrderItems, customError *exceptions.CustomError) {
	customError = i.IsOrderExist(ctx, orderID)
	if customError != nil {
		return
	}

	items, err := i.repository.GetOrderWithItems(ctx, orderID)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get order items: %w", err),
		}
	}

	result = make([]*domain.OrderItems, len(items))
	for index, v := range items {
		result[index] = &domain.OrderItems{
			ID:            v.ID,
			OrderID:       v.OrderID,
			ProductID:     v.ProductID,
			StatusID:      v.StatusID,
			TableNumber:   tableNumber,
			StatusName:    v.StatusName,
			StatusNameEN:  v.StatusNameEN,
			StatusCode:    v.StatusCode,
			ProductName:   v.ProductName,
			ProductNameEN: v.ProductNameEN,
			Price:         utils.PgNumericToFloat64(v.Price),
			Quantity:      v.Quantity,
			Note:          utils.PgTextToStringPtr(v.Note),
		}
	}

	return
}

func (i *Implement) GetOderItemsByID(ctx context.Context, orderID, orderItemsID int64, tableNumber int32) (result *domain.OrderItems, customError *exceptions.CustomError) {

	customError = i.IsOrderExist(ctx, orderID)
	if customError != nil {
		return
	}

	items, err := i.repository.GetOrderWithItemsByID(ctx, database.GetOrderWithItemsByIDParams{
		OrderID:      orderID,
		OrderItemsID: orderItemsID,
	})
	if err != nil {

		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: fmt.Errorf("order items not found"),
			}
		}

		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get order items: %w", err),
		}
	}

	return &domain.OrderItems{
		ID:            items.ID,
		OrderID:       items.OrderID,
		ProductID:     items.ProductID,
		StatusID:      items.StatusID,
		TableNumber:   tableNumber,
		StatusName:    items.StatusName,
		StatusNameEN:  items.StatusNameEN,
		StatusCode:    items.StatusCode,
		ProductName:   items.ProductName,
		ProductNameEN: items.ProductNameEN,
		Price:         utils.PgNumericToFloat64(items.Price),
		Quantity:      items.Quantity,
		Note:          utils.PgTextToStringPtr(items.Note),
	}, nil
}

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload domain.OrderItemsStatus) (customError *exceptions.CustomError) {
	customError = i.IsOrderExist(ctx, payload.OrderID)
	if customError != nil {
		return
	}

	customError = i.IsOrderStatus(ctx, payload.StatusCode)
	if customError != nil {
		return
	}

	err := i.repository.UpdateOrderItemsStatus(ctx, database.UpdateOrderItemsStatusParams{
		StatusCode: payload.StatusCode,
		ID:         payload.ID,
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update order items status: %w", err),
		}
	}

	return
}
