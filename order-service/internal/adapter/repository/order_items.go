package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	shareModel "food-story/shared/model"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

const FailedToGetOrderItems = "failed to get order items: %w"

func (i *Implement) CreateOrderItems(ctx context.Context, orderItems []shareModel.OrderItems) (result []*shareModel.OrderItems, customError *exceptions.CustomError) {

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

	tableID, err := i.repository.GetTableIDByOrderID(ctx, orderItemsPayload[0].OrderID)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get table id by order id: %w", err),
		}
	}

	err = i.repository.UpdateTablesStatusWaitingToBeServed(ctx, tableID)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update tables status waiting to be served: %w", err),
		}
	}

	var orderItemsID []int64
	for _, item := range orderItemsPayload {
		orderItemsID = append(orderItemsID, item.ID)
	}

	result, customError = i.GetOderItemsGroupID(ctx, orderItemsID)
	if customError != nil {
		return nil, customError
	}

	return
}

func (i *Implement) GetOrderItemsByOrderID(ctx context.Context, orderID int64) (result []*shareModel.OrderItems, customError *exceptions.CustomError) {

	customError = i.validateAndCheckOrder(ctx, orderID)
	if customError != nil {
		return nil, customError
	}

	orderItems, repoErr := i.repository.GetOrderItemsByOrderID(ctx, orderID)
	if repoErr != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf(FailedToGetOrderItems, repoErr),
		}
	}

	return shareModel.TransformOrderItemsResults(orderItems), nil
}

func (i *Implement) GetOderItemsGroupID(ctx context.Context, orderItemsID []int64) (result []*shareModel.OrderItems, customError *exceptions.CustomError) {

	if len(orderItemsID) == 0 {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: fmt.Errorf("order items id cannot be empty"),
		}
	}

	orderItems, repoErr := i.repository.GetOrderWithItemsGroupID(ctx, orderItemsID)
	if repoErr != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf(FailedToGetOrderItems, repoErr),
		}
	}

	return shareModel.TransformOrderItemsResults(orderItems), nil
}

func (i *Implement) GetCurrentOrderItems(ctx context.Context, orderID int64, pageNumberParam, pageSizeParam int64) (result domain.SearchCurrentOrderItemsResult, customError *exceptions.CustomError) {
	pageSize, pageNumber := utils.CalculatePageSizeAndNumber(pageSizeParam, pageNumberParam)

	customError = i.validateAndCheckOrder(ctx, orderID)
	if customError != nil {
		return domain.SearchCurrentOrderItemsResult{}, customError
	}

	var (
		searchResult []*database.GetOrderWithItemsRow
		totalItems   int64
	)

	// parallel search
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		searchResult, customError = i.fetchOrderWithItems(ctx, database.GetOrderWithItemsParams{
			OrderID:    orderID,
			Pagesize:   pageSize,
			Pagenumber: pageNumber,
		})
	}()
	go func() {
		defer wg.Done()
		totalItems, customError = i.fetchTotalOrderWithItems(ctx, orderID)
	}()

	wg.Wait()

	if customError != nil {
		return domain.SearchCurrentOrderItemsResult{}, customError
	}

	return domain.SearchCurrentOrderItemsResult{
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, pageSize),
		Data:       transformOrderItemsResults(searchResult),
	}, nil
}

func (i *Implement) GetCurrentOrderItemsByID(ctx context.Context, orderID, orderItemsID int64) (result *domain.CurrentOrderItems, customError *exceptions.CustomError) {

	orderItem, repoErr := i.repository.GetOrderWithItemsByID(ctx, database.GetOrderWithItemsByIDParams{
		OrderID:      orderID,
		OrderItemsID: orderItemsID,
	})
	if repoErr != nil {
		if errors.Is(repoErr, exceptions.ErrRowDatabaseNotFound) {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrOrderItemsNotFound,
			}
		}

		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf(FailedToGetOrderItems, repoErr),
		}
	}

	return transformCurrentOrderItemsByIDResults(orderItem), nil
}

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload shareModel.OrderItemsStatus) (customError *exceptions.CustomError) {
	customError = i.IsOrderWithItemsExists(ctx, payload.OrderID, payload.ID)
	if customError != nil {
		return customError
	}

	repoErr := i.repository.UpdateOrderItemsStatus(ctx, database.UpdateOrderItemsStatusParams{
		StatusCode: payload.StatusCode,
		ID:         payload.ID,
	})
	if repoErr != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update order items status: %w", repoErr),
		}
	}

	return
}

func (i *Implement) SearchOrderItemsIncomplete(ctx context.Context, orderID int64, search domain.SearchOrderItems) (result domain.SearchOrderItemsResult, customError *exceptions.CustomError) {
	searchParams := buildSearchOrderItemsParams(orderID, search)

	var (
		searchResult []*database.SearchOrderItemsIsNotFinalRow
		totalItems   int64
	)

	// parallel search
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		searchResult, customError = i.fetchOrderItemsNotFinal(ctx, searchParams)
	}()
	go func() {
		defer wg.Done()
		totalItems, customError = i.fetchTotalItemsNotFinal(ctx, searchParams)
	}()

	wg.Wait()

	if customError != nil {
		return domain.SearchOrderItemsResult{}, customError
	}

	return domain.SearchOrderItemsResult{
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       shareModel.TransformOrderItemsResults(searchResult),
	}, nil
}

func (i *Implement) fetchOrderWithItems(ctx context.Context, params database.GetOrderWithItemsParams) ([]*database.GetOrderWithItemsRow, *exceptions.CustomError) {
	result, err := i.repository.GetOrderWithItems(ctx, params)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch order items not final: %w", err),
		}
	}
	return result, nil
}

func (i *Implement) fetchTotalOrderWithItems(ctx context.Context, orderID int64) (int64, *exceptions.CustomError) {
	totalItems, err := i.repository.GetTotalItemOrderWithItems(ctx, orderID)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch total items: %w", err),
		}
	}

	return totalItems, nil
}

func (i *Implement) fetchOrderItemsNotFinal(ctx context.Context, params database.SearchOrderItemsIsNotFinalParams) ([]*database.SearchOrderItemsIsNotFinalRow, *exceptions.CustomError) {
	result, err := i.repository.SearchOrderItemsIsNotFinal(ctx, params)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch order items not final: %w", err),
		}
	}
	return result, nil
}

func (i *Implement) fetchTotalItemsNotFinal(ctx context.Context, params database.SearchOrderItemsIsNotFinalParams) (int64, *exceptions.CustomError) {
	totalParams := database.GetTotalSearchOrderItemsIsNotFinalParams{
		ProductName: params.ProductName,
		OrderID:     params.OrderID,
		StatusCode:  params.StatusCode,
	}
	totalItems, err := i.repository.GetTotalSearchOrderItemsIsNotFinal(ctx, totalParams)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch total items: %w", err),
		}
	}

	return totalItems, nil
}

func (i *Implement) buildPayloadOrderItems(ctx context.Context, orderItems []shareModel.OrderItems) ([]database.CreateOrderItemsParams, *exceptions.CustomError) {

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

	result := make([]database.CreateOrderItemsParams, len(orderItems))
	for index, item := range orderItems {
		product, repoErr := i.repository.GetProductByID(ctx, item.ProductID)
		if repoErr != nil {
			if errors.Is(repoErr, exceptions.ErrRowDatabaseNotFound) {
				return []database.CreateOrderItemsParams{}, &exceptions.CustomError{
					Status: exceptions.ERRNOTFOUND,
					Errors: exceptions.ErrProductNotFound,
				}
			}
			return []database.CreateOrderItemsParams{}, &exceptions.CustomError{
				Status: exceptions.ERRREPOSITORY,
				Errors: fmt.Errorf("failed to get product by id: %w", repoErr),
			}
		}

		if product == nil {
			return []database.CreateOrderItemsParams{}, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrProductNotFound,
			}
		}

		result[index] = database.CreateOrderItemsParams{
			ID:              i.snowflakeID.Generate(),
			OrderID:         item.OrderID,
			ProductID:       product.ID,
			StatusID:        statusPreparingID,
			ProductName:     product.Name,
			ProductNameEn:   product.NameEn,
			Price:           product.Price,
			Quantity:        item.Quantity,
			Note:            utils.StringPtrToPgText(item.Note),
			CreatedAt:       currentTime,
			ProductImageUrl: product.ImageUrl,
			IsVisible:       true,
		}
	}

	return result, nil
}

func buildSearchOrderItemsParams(orderID int64, payload domain.SearchOrderItems) database.SearchOrderItemsIsNotFinalParams {
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

func validationOrderItems(items []shareModel.OrderItems) *exceptions.CustomError {
	if len(items) == 0 {
		return &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: exceptions.ErrOrderItemsRequired,
		}
	}
	return nil
}

func (i *Implement) validateAndCheckOrder(ctx context.Context, orderID int64) *exceptions.CustomError {
	if orderID <= 0 {
		return &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: exceptions.ErrOrderRequired,
		}
	}

	return i.IsOrderExist(ctx, orderID)
}

type CurrentOrderItemsRow interface {
	GetID() int64
	GetProductID() int64
	GetStatusName() string
	GetStatusNameEN() string
	GetStatusCode() string
	GetProductName() string
	GetProductNameEN() string
	GetImageURL() pgtype.Text
	GetPrice() pgtype.Numeric
	GetQuantity() int32
	GetNote() pgtype.Text
	GetCreatedAt() pgtype.Timestamptz
}

func transformCurrentOrderItemsByIDResults[T CurrentOrderItemsRow](results T) *domain.CurrentOrderItems {
	createdAt, _ := utils.PgTimestampToThaiISO8601(results.GetCreatedAt())
	return &domain.CurrentOrderItems{
		ID:            results.GetID(),
		ProductID:     results.GetProductID(),
		StatusName:    results.GetStatusName(),
		StatusNameEN:  results.GetStatusNameEN(),
		StatusCode:    results.GetStatusCode(),
		ProductName:   results.GetProductName(),
		ProductNameEN: results.GetProductNameEN(),
		ImageURL:      utils.PgTextToStringPtr(results.GetImageURL()),
		Price:         utils.PgNumericToFloat64(results.GetPrice()),
		Quantity:      results.GetQuantity(),
		Note:          utils.PgTextToStringPtr(results.GetNote()),
		CreatedAt:     createdAt,
	}
}

func transformOrderItemsResults[T CurrentOrderItemsRow](results []T) []*domain.CurrentOrderItems {
	data := make([]*domain.CurrentOrderItems, len(results))
	for index, row := range results {
		data[index] = transformCurrentOrderItemsByIDResults(row)
	}
	return data
}
