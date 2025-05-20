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
)

func (i *Implement) CreateOrderItems(ctx context.Context, orderItems []domain.OrderItems, tableNumber int32) (result []*domain.OrderItems, customError *exceptions.CustomError) {

	validationError := validationOrderItems(orderItems)
	if validationError != nil {
		return nil, validationError
	}

	orderItemsPayload, buildParamError := i.buildPayloadOrderItems(ctx, orderItems)
	if buildParamError != nil {
		return nil, buildParamError
	}

	_, err := i.repository.CreateOrderItems(ctx, orderItemsPayload)
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

func (i *Implement) buildPayloadOrderItems(ctx context.Context, orderItems []domain.OrderItems) ([]database.CreateOrderItemsParams, *exceptions.CustomError) {

	validationError := validationOrderItems(orderItems)
	if validationError != nil {
		return []database.CreateOrderItemsParams{}, validationError
	}

	statusPreparingID, statusIDErr := i.GetOrderStatusPreparing(ctx)
	if statusIDErr != nil {
		return []database.CreateOrderItemsParams{}, statusIDErr
	}

	currentTime, timeErr := i.repository.GetTimeNow(ctx)
	if timeErr != nil {
		return []database.CreateOrderItemsParams{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get current time: %w", timeErr),
		}
	}

	result := make([]database.CreateOrderItemsParams, 0, len(orderItems))
	for index, item := range orderItems {
		product, repoErr := i.repository.GetProductByID(ctx, item.ProductID)
		if repoErr != nil || product == nil {
			msg := fmt.Sprintf("product %d not found", item.ProductID)
			status := exceptions.ERRNOTFOUND

			if repoErr != nil && !errors.Is(repoErr, exceptions.ErrRowDatabaseNotFound) {
				status = exceptions.ERRREPOSITORY
				msg = fmt.Sprintf("failed to get product: %v", repoErr)
			}

			return []database.CreateOrderItemsParams{}, &exceptions.CustomError{
				Status: status,
				Errors: fmt.Errorf(msg),
			}
		}

		result[index] = database.CreateOrderItemsParams{
			ID:            i.snowflakeID.Generate(),
			OrderID:       item.OrderID,
			ProductID:     product.ID,
			StatusID:      statusPreparingID,
			ProductName:   product.Name,
			ProductNameEn: product.NameEn,
			Price:         product.Price,
			Quantity:      item.Quantity,
			Note:          utils.StringPtrToPgText(item.Note),
			CreatedAt:     currentTime,
		}
	}

	return result, nil
}

func (i *Implement) GetOrderItems(ctx context.Context, orderID int64, tableNumber int32) (result []*domain.OrderItems, customError *exceptions.CustomError) {

	if orderID <= 0 || tableNumber <= 0 {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: fmt.Errorf("order id or table number cannot be empty"),
		}
	}

	customError = i.IsOrderExist(ctx, orderID)
	if customError != nil {
		return nil, customError
	}

	orderItems, repoErr := i.repository.GetOrderWithItems(ctx, orderID)
	if repoErr != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get order items: %w", repoErr),
		}
	}

	result = make([]*domain.OrderItems, len(orderItems))
	for index, item := range orderItems {
		createdAt, timeErr := utils.PgTimestampToThaiISO8601(item.CreatedAt)
		if timeErr != nil {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRSYSTEM,
				Errors: timeErr,
			}
		}

		result[index] = &domain.OrderItems{
			ID:            item.ID,
			OrderID:       item.OrderID,
			OrderNumber:   item.OrderNumber,
			ProductID:     item.ProductID,
			StatusID:      item.StatusID,
			TableNumber:   tableNumber,
			StatusName:    item.StatusName,
			StatusNameEN:  item.StatusNameEN,
			StatusCode:    item.StatusCode,
			ProductName:   item.ProductName,
			ProductNameEN: item.ProductNameEN,
			Price:         utils.PgNumericToFloat64(item.Price),
			Quantity:      item.Quantity,
			Note:          utils.PgTextToStringPtr(item.Note),
			CreatedAt:     createdAt,
		}

	}

	return
}

func (i *Implement) GetCurrentOrderItems(ctx context.Context, orderID int64) (result []*domain.CurrentOrderItems, customError *exceptions.CustomError) {

	if orderID <= 0 {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: fmt.Errorf("order id cannot be empty"),
		}
	}

	customError = i.IsOrderExist(ctx, orderID)
	if customError != nil {
		return nil, customError
	}

	orderItems, repoErr := i.repository.GetOrderWithItems(ctx, orderID)
	if repoErr != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get order items: %w", repoErr),
		}
	}

	result = make([]*domain.CurrentOrderItems, len(orderItems))
	for index, item := range orderItems {
		createdAt, timErr := utils.PgTimestampToThaiISO8601(item.CreatedAt)
		if timErr != nil {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRSYSTEM,
				Errors: timErr,
			}
		}

		result[index] = &domain.CurrentOrderItems{
			ID:            item.ID,
			ProductID:     item.ProductID,
			StatusName:    item.StatusName,
			StatusNameEN:  item.StatusNameEN,
			StatusCode:    item.StatusCode,
			ProductName:   item.ProductName,
			ProductNameEN: item.ProductNameEN,
			Price:         utils.PgNumericToFloat64(item.Price),
			Quantity:      item.Quantity,
			Note:          utils.PgTextToStringPtr(item.Note),
			CreatedAt:     createdAt,
		}

	}

	return
}

func (i *Implement) GetOderItemsByID(ctx context.Context, orderID, orderItemsID int64) (result *domain.CurrentOrderItems, customError *exceptions.CustomError) {

	orderItemsExistsErr := i.IsOrderWithItemsExists(ctx, orderID, orderItemsID)
	if orderItemsExistsErr != nil {
		return nil, orderItemsExistsErr
	}

	orderItem, repoErr := i.repository.GetOrderWithItemsByID(ctx, database.GetOrderWithItemsByIDParams{
		OrderID:      orderID,
		OrderItemsID: orderItemsID,
	})
	if repoErr != nil || orderItem == nil {
		msg := fmt.Sprintf("order items is nil")
		status := exceptions.ERRNOTFOUND

		if repoErr != nil && !errors.Is(repoErr, exceptions.ErrRowDatabaseNotFound) {
			status = exceptions.ERRREPOSITORY
			msg = fmt.Sprintf("failed to get order items: %v", repoErr)
		}

		return nil, &exceptions.CustomError{
			Status: status,
			Errors: fmt.Errorf(msg),
		}
	}

	createdAt, timeErr := utils.PgTimestampToThaiISO8601(orderItem.CreatedAt)
	if timeErr != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: timeErr,
		}
	}

	return &domain.CurrentOrderItems{
		ID:            orderItem.ID,
		ProductID:     orderItem.ProductID,
		StatusName:    orderItem.StatusName,
		StatusNameEN:  orderItem.StatusNameEN,
		StatusCode:    orderItem.StatusCode,
		ProductName:   orderItem.ProductName,
		ProductNameEN: orderItem.ProductNameEN,
		Price:         utils.PgNumericToFloat64(orderItem.Price),
		Quantity:      orderItem.Quantity,
		Note:          utils.PgTextToStringPtr(orderItem.Note),
		CreatedAt:     createdAt,
	}, nil
}

func (i *Implement) GetOderItemsGroupID(ctx context.Context, orderItemsID []int64, tableNumber int32) (result []*domain.OrderItems, customError *exceptions.CustomError) {

	orderItems, repoErr := i.repository.GetOrderWithItemsGroupID(ctx, orderItemsID)
	if repoErr != nil || orderItems == nil {
		msg := fmt.Sprintf("order items is nil")
		status := exceptions.ERRNOTFOUND

		if repoErr != nil && !errors.Is(repoErr, exceptions.ErrRowDatabaseNotFound) {
			status = exceptions.ERRREPOSITORY
			msg = fmt.Sprintf("failed to get order items: %v", repoErr)
		}

		return nil, &exceptions.CustomError{
			Status: status,
			Errors: fmt.Errorf(msg),
		}
	}

	result = make([]*domain.OrderItems, len(orderItems))
	for index, item := range orderItems {
		createdAt, err := utils.PgTimestampToThaiISO8601(item.CreatedAt)
		if err != nil {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRSYSTEM,
				Errors: err,
			}
		}

		result[index] = &domain.OrderItems{
			ID:            item.ID,
			OrderID:       item.OrderID,
			OrderNumber:   item.OrderNumber,
			ProductID:     item.ProductID,
			StatusID:      item.StatusID,
			TableNumber:   tableNumber,
			StatusName:    item.StatusName,
			StatusNameEN:  item.StatusNameEN,
			StatusCode:    item.StatusCode,
			ProductName:   item.ProductName,
			ProductNameEN: item.ProductNameEN,
			Price:         utils.PgNumericToFloat64(item.Price),
			Quantity:      item.Quantity,
			Note:          utils.PgTextToStringPtr(item.Note),
			CreatedAt:     createdAt,
		}
	}
	return
}

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload domain.OrderItemsStatus) (customError *exceptions.CustomError) {
	customError = i.IsOrderWithItemsExists(ctx, payload.OrderID, payload.ID)
	if customError != nil {
		return customError
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

func validationOrderItems(items []domain.OrderItems) *exceptions.CustomError {
	if len(items) == 0 {
		return &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: exceptions.ErrOrderItemsRequired,
		}
	}
	return nil
}
