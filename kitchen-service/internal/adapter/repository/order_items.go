package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
	"strings"
)

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

func (i *Implement) UpdateOrderItemsStatusServed(ctx context.Context, payload domain.OrderItemsStatus) (customError *exceptions.CustomError) {
	customError = i.IsOrderWithItemsExists(ctx, payload.OrderID, payload.ID)
	if customError != nil {
		return
	}

	err := i.repository.UpdateOrderItemsStatusServed(ctx, payload.ID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update order items status: %w", err),
		}
	}

	return
}

func (i *Implement) SearchOrderItems(ctx context.Context, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, customError *exceptions.CustomError) {
	searchParams := buildSearchOrderItemsParams(payload)

	items, err := i.repository.SearchOrderItems(ctx, searchParams)
	if err != nil {
		return domain.SearchOrderItemsResult{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get order items: %w", err),
		}
	}

	totalItemsParam := database.GetTotalSearchOrderItemsParams{
		ProductName: searchParams.ProductName,
		TableNumber: searchParams.TableNumber,
		StatusCode:  searchParams.StatusCode,
	}

	totalItems, err := i.repository.GetTotalSearchOrderItems(ctx, totalItemsParam)
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
				Status: exceptions.ERRUNKNOWN,
				Errors: err,
			}
		}

		data[index] = &domain.OrderItems{
			ID:            v.ID,
			OrderID:       v.OrderID,
			OrderNumber:   v.OrderNumber,
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
				Status: exceptions.ERRUNKNOWN,
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

func (i *Implement) GetOrderItemsByID(ctx context.Context, orderID, orderItemsID int64, tableNumber int32) (result *domain.OrderItems, customError *exceptions.CustomError) {

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
			Status: exceptions.ERRUNKNOWN,
			Errors: err,
		}
	}

	return &domain.OrderItems{
		ID:            items.ID,
		OrderID:       items.OrderID,
		OrderNumber:   items.OrderNumber,
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
		CreatedAt:     createdAt,
	}, nil
}

func buildSearchOrderItemsParams(payload domain.SearchOrderItems) database.SearchOrderItemsParams {
	params := database.SearchOrderItemsParams{
		ProductName: pgtype.Text{String: payload.Name, Valid: payload.Name != ""},
		TableNumber: payload.TableNumber,
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
