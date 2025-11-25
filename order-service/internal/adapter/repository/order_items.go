package repository

import (
	"context"
	"errors"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	shareModel "food-story/shared/model"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

const FailedToGetOrderItems = "failed to get order items"

func (i *Implement) CreateOrderItems(ctx context.Context, orderItems []shareModel.OrderItems) (result []*shareModel.OrderItems, err error) {

	err = validationOrderItems(orderItems)
	if err != nil {
		return nil, err
	}

	orderItemsPayload, buildParamError := i.buildPayloadOrderItems(ctx, orderItems)
	if buildParamError != nil {
		return nil, buildParamError
	}

	_, err = i.repository.CreateOrderItems(ctx, orderItemsPayload)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to create order items", err)
	}

	tableID, err := i.repository.GetTableIDByOrderID(ctx, orderItemsPayload[0].OrderID)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to get table id by order id", err)
	}

	err = i.repository.UpdateTablesStatusWaitingToBeServed(ctx, tableID)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to update tables status waiting to be served", err)
	}

	var orderItemsID []int64
	for _, item := range orderItemsPayload {
		orderItemsID = append(orderItemsID, item.ID)
	}

	result, err = i.GetOderItemsGroupID(ctx, orderItemsID)
	if err != nil {
		return nil, err
	}

	return
}

func (i *Implement) GetOrderItemsByOrderID(ctx context.Context, orderID int64) (result []*shareModel.OrderItems, err error) {

	err = i.validateAndCheckOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	orderItems, err := i.repository.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, FailedToGetOrderItems, err)
	}

	return shareModel.TransformOrderItemsResults(orderItems), nil
}

func (i *Implement) GetOderItemsGroupID(ctx context.Context, orderItemsID []int64) (result []*shareModel.OrderItems, err error) {

	if len(orderItemsID) == 0 {
		return nil, exceptions.Error(exceptions.CodeBusiness, "order items id cannot be empty")
	}

	orderItems, err := i.repository.GetOrderWithItemsGroupID(ctx, orderItemsID)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, FailedToGetOrderItems, err)
	}

	return shareModel.TransformOrderItemsResults(orderItems), nil
}

func (i *Implement) GetCurrentOrderItems(ctx context.Context, orderID int64, pageNumberParam, pageSizeParam int64) (result domain.SearchCurrentOrderItemsResult, err error) {
	pageSize, pageNumber := utils.CalculatePageSizeAndNumber(pageSizeParam, pageNumberParam)

	err = i.validateAndCheckOrder(ctx, orderID)
	if err != nil {
		return domain.SearchCurrentOrderItemsResult{}, err
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
		searchResult, err = i.fetchOrderWithItems(ctx, database.GetOrderWithItemsParams{
			OrderID:    orderID,
			Pagesize:   pageSize,
			Pagenumber: pageNumber,
		})
	}()
	go func() {
		defer wg.Done()
		totalItems, err = i.fetchTotalOrderWithItems(ctx, orderID)
	}()

	wg.Wait()

	if err != nil {
		return domain.SearchCurrentOrderItemsResult{}, err
	}

	return domain.SearchCurrentOrderItemsResult{
		PageNumber: utils.GetPageNumber(pageNumberParam),
		PageSize:   utils.GetPageSize(pageSizeParam),
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, pageSize),
		Data:       transformOrderItemsResults(searchResult),
	}, nil
}

func (i *Implement) GetCurrentOrderItemsByID(ctx context.Context, orderID, orderItemsID int64) (result *domain.CurrentOrderItems, err error) {

	orderItem, err := i.repository.GetOrderWithItemsByID(ctx, database.GetOrderWithItemsByIDParams{
		OrderID:      orderID,
		OrderItemsID: orderItemsID,
	})
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderItemsNotFound.Error())
		}

		return nil, exceptions.Errorf(exceptions.CodeRepository, FailedToGetOrderItems, err)
	}

	return transformCurrentOrderItemsByIDResults(orderItem), nil
}

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload shareModel.OrderItemsStatus) (err error) {
	err = i.IsOrderWithItemsExists(ctx, payload.OrderID, payload.ID)
	if err != nil {
		return err
	}

	err = i.repository.UpdateOrderItemsStatus(ctx, database.UpdateOrderItemsStatusParams{
		StatusCode: payload.StatusCode,
		ID:         payload.ID,
	})
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to update order items status", err)
	}

	return
}

func (i *Implement) SearchOrderItemsIncomplete(ctx context.Context, orderID int64, search domain.SearchOrderItems) (result domain.SearchOrderItemsResult, err error) {
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
		searchResult, err = i.fetchOrderItemsNotFinal(ctx, searchParams)
	}()
	go func() {
		defer wg.Done()
		totalItems, err = i.fetchTotalItemsNotFinal(ctx, searchParams)
	}()

	wg.Wait()

	if err != nil {
		return domain.SearchOrderItemsResult{}, err
	}

	return domain.SearchOrderItemsResult{
		PageNumber: utils.GetPageNumber(search.PageNumber),
		PageSize:   utils.GetPageSize(search.PageSize),
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       shareModel.TransformOrderItemsResults(searchResult),
	}, nil
}

func (i *Implement) fetchOrderWithItems(ctx context.Context, params database.GetOrderWithItemsParams) ([]*database.GetOrderWithItemsRow, error) {
	result, err := i.repository.GetOrderWithItems(ctx, params)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch order items not final", err)
	}
	return result, nil
}

func (i *Implement) fetchTotalOrderWithItems(ctx context.Context, orderID int64) (int64, error) {
	totalItems, err := i.repository.GetTotalItemOrderWithItems(ctx, orderID)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch total items", err)
	}

	return totalItems, nil
}

func (i *Implement) fetchOrderItemsNotFinal(ctx context.Context, params database.SearchOrderItemsIsNotFinalParams) ([]*database.SearchOrderItemsIsNotFinalRow, error) {
	result, err := i.repository.SearchOrderItemsIsNotFinal(ctx, params)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch order items not final", err)
	}
	return result, nil
}

func (i *Implement) fetchTotalItemsNotFinal(ctx context.Context, params database.SearchOrderItemsIsNotFinalParams) (int64, error) {
	totalParams := database.GetTotalSearchOrderItemsIsNotFinalParams{
		ProductName: params.ProductName,
		OrderID:     params.OrderID,
		StatusCode:  params.StatusCode,
	}
	totalItems, err := i.repository.GetTotalSearchOrderItemsIsNotFinal(ctx, totalParams)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch total items order not final", err)
	}

	return totalItems, nil
}

func (i *Implement) buildPayloadOrderItems(ctx context.Context, orderItems []shareModel.OrderItems) ([]database.CreateOrderItemsParams, error) {

	err := validationOrderItems(orderItems)
	if err != nil {
		return []database.CreateOrderItemsParams{}, err
	}

	statusPreparingID, err := i.GetOrderStatusPreparing(ctx)
	if err != nil {
		return []database.CreateOrderItemsParams{}, err
	}

	currentTime, timeErr := i.repository.GetTimeNow(ctx)
	if timeErr != nil {
		return []database.CreateOrderItemsParams{}, exceptions.Errorf(exceptions.CodeRepository, "failed to get current time", err)
	}

	result := make([]database.CreateOrderItemsParams, len(orderItems))
	for index, item := range orderItems {
		product, err := i.repository.GetProductByID(ctx, item.ProductID)
		if err != nil {
			if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
				return []database.CreateOrderItemsParams{}, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrProductNotFound.Error())
			}

			return []database.CreateOrderItemsParams{}, exceptions.Errorf(exceptions.CodeRepository, "failed to get product by id", err)
		}

		if product == nil {
			return []database.CreateOrderItemsParams{}, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrProductNotFound.Error())
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

func validationOrderItems(items []shareModel.OrderItems) error {
	if len(items) == 0 {
		return exceptions.Error(exceptions.CodeBusiness, exceptions.ErrOrderItemsRequired.Error())
	}
	return nil
}

func (i *Implement) validateAndCheckOrder(ctx context.Context, orderID int64) error {
	if orderID <= 0 {
		return exceptions.Error(exceptions.CodeBusiness, exceptions.ErrOrderRequired.Error())
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
	if len(data) == 0 {
		return []*domain.CurrentOrderItems{}
	}
	return data
}
