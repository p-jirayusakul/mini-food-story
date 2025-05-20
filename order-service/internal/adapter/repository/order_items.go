package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
	"strings"
	"time"
)

func (i *Implement) CreateOrderItems(ctx context.Context, items []domain.OrderItems, tableNumber int32) (result []*domain.OrderItems, customError *exceptions.CustomError) {

	if len(items) == 0 {
		return nil, nil
	}

	orderItems, buildParamError := i.buildPayloadOrderItems(ctx, items)
	if buildParamError != nil {
		return nil, buildParamError
	}

	_, err := i.repository.CreateOrderItems(ctx, orderItems)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to create order items: %w", err),
		}
	}

	var orderItemsID []int64
	for _, item := range orderItems {
		orderItemsID = append(orderItemsID, item.ID)
	}

	result, customError = i.GetOderItemsGroupID(ctx, orderItemsID, tableNumber)
	if customError != nil {
		return nil, customError
	}

	return
}

func (i *Implement) buildPayloadOrderItems(ctx context.Context, items []domain.OrderItems) ([]database.CreateOrderItemsParams, *exceptions.CustomError) {
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
			CreatedAt:     pgtype.Timestamptz{Time: time.Now(), Valid: true},
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
		createdAt, err := utils.PgTimestampToThaiISO8601(v.CreatedAt)
		if err != nil {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRSYSTEM,
				Errors: err,
			}
		}

		result[index] = &domain.OrderItems{
			ID:            v.ID,
			OrderID:       v.OrderID,
			OrderNumber:   v.OrderNumber,
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
			CreatedAt:     createdAt,
		}

	}

	return
}

func (i *Implement) GetCurrentOrderItems(ctx context.Context, orderID int64, tableNumber int32) (result []*domain.CurrentOrderItems, customError *exceptions.CustomError) {
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

	result = make([]*domain.CurrentOrderItems, len(items))
	for index, v := range items {
		createdAt, err := utils.PgTimestampToThaiISO8601(v.CreatedAt)
		if err != nil {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRSYSTEM,
				Errors: err,
			}
		}

		result[index] = &domain.CurrentOrderItems{
			ID:            v.ID,
			ProductID:     v.ProductID,
			StatusName:    v.StatusName,
			StatusNameEN:  v.StatusNameEN,
			StatusCode:    v.StatusCode,
			ProductName:   v.ProductName,
			ProductNameEN: v.ProductNameEN,
			Price:         utils.PgNumericToFloat64(v.Price),
			Quantity:      v.Quantity,
			Note:          utils.PgTextToStringPtr(v.Note),
			CreatedAt:     createdAt,
		}

	}

	return
}

func (i *Implement) GetOderItemsByID(ctx context.Context, orderID, orderItemsID int64, tableNumber int32) (result *domain.CurrentOrderItems, customError *exceptions.CustomError) {

	customError = i.IsOrderWithItemsExists(ctx, orderID, orderItemsID)
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

	createdAt, err := utils.PgTimestampToThaiISO8601(items.CreatedAt)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: err,
		}
	}

	return &domain.CurrentOrderItems{
		ID:            items.ID,
		ProductID:     items.ProductID,
		StatusName:    items.StatusName,
		StatusNameEN:  items.StatusNameEN,
		StatusCode:    items.StatusCode,
		ProductName:   items.ProductName,
		ProductNameEN: items.ProductNameEN,
		Price:         utils.PgNumericToFloat64(items.Price),
		Quantity:      items.Quantity,
		Note:          utils.PgTextToStringPtr(items.Note),
		CreatedAt:     createdAt,
	}, nil
}

func (i *Implement) GetOderItemsGroupID(ctx context.Context, orderItemsID []int64, tableNumber int32) (result []*domain.OrderItems, customError *exceptions.CustomError) {

	items, err := i.repository.GetOrderWithItemsGroupID(ctx, orderItemsID)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get order items: %w", err),
		}
	}

	result = make([]*domain.OrderItems, len(items))
	for index, v := range items {
		createdAt, err := utils.PgTimestampToThaiISO8601(v.CreatedAt)
		if err != nil {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRSYSTEM,
				Errors: err,
			}
		}

		result[index] = &domain.OrderItems{
			ID:            v.ID,
			OrderID:       v.OrderID,
			OrderNumber:   v.OrderNumber,
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
			CreatedAt:     createdAt,
		}
	}
	return
}

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload domain.OrderItemsStatus) (customError *exceptions.CustomError) {
	customError = i.IsOrderWithItemsExists(ctx, payload.OrderID, payload.ID)
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

func (i *Implement) SearchOrderItemsIncomplete(ctx context.Context, orderID int64, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, customError *exceptions.CustomError) {
	searchParams := buildSearchOrderItemsIncompleteParams(orderID, payload)

	items, err := i.repository.SearchOrderItemsIsNotFinal(ctx, searchParams)
	if err != nil {
		return domain.SearchOrderItemsResult{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get order items: %w", err),
		}
	}

	totalItemsParam := database.GetTotalSearchOrderItemsIsNotFinalParams{
		ProductName: searchParams.ProductName,
		OrderID:     orderID,
		StatusCode:  searchParams.StatusCode,
	}

	totalItems, err := i.repository.GetTotalSearchOrderItemsIsNotFinal(ctx, totalItemsParam)
	if err != nil {
		return domain.SearchOrderItemsResult{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch product: %w", err),
		}
	}

	data := make([]*domain.OrderItems, len(items))
	for index, v := range items {
		createdAt, err := utils.PgTimestampToThaiISO8601(v.CreatedAt)
		if err != nil {
			return domain.SearchOrderItemsResult{}, &exceptions.CustomError{
				Status: exceptions.ERRSYSTEM,
				Errors: err,
			}
		}

		data[index] = &domain.OrderItems{
			ID:            v.ID,
			OrderID:       v.OrderID,
			ProductID:     v.ProductID,
			StatusID:      v.StatusID,
			TableNumber:   v.TableNumber,
			StatusName:    v.StatusName,
			StatusNameEN:  v.StatusNameEN,
			StatusCode:    v.StatusCode,
			ProductName:   v.ProductName,
			ProductNameEN: v.ProductNameEN,
			Price:         utils.PgNumericToFloat64(v.Price),
			Quantity:      v.Quantity,
			Note:          utils.PgTextToStringPtr(v.Note),
			CreatedAt:     createdAt,
		}

	}

	return domain.SearchOrderItemsResult{
		TotalItems: totalItems,
		TotalPages: int64(math.Ceil(float64(totalItems) / float64(searchParams.PageSize))),
		Data:       data,
	}, nil
}

func buildSearchOrderItemsIncompleteParams(orderID int64, payload domain.SearchOrderItems) database.SearchOrderItemsIsNotFinalParams {
	params := database.SearchOrderItemsIsNotFinalParams{
		OrderID:     orderID,
		ProductName: pgtype.Text{String: payload.Name, Valid: payload.Name != ""},
		OrderByType: payload.OrderByType,
		OrderBy:     payload.OrderBy,
		PageSize:    payload.PageSize,
		PageNumber:  payload.PageNumber,
	}

	for _, v := range payload.StatusCode {
		params.StatusCode = append(params.StatusCode, strings.ToUpper(v))
	}

	params.PageSize, params.PageNumber = utils.CalculatePageSizeAndNumber(payload.PageSize, payload.PageNumber)

	return params
}
