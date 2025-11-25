package repository

import (
	"context"
	"errors"
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	shareModel "food-story/shared/model"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload shareModel.OrderItemsStatus) (err error) {
	err = i.IsOrderWithItemsExists(ctx, payload.OrderID, payload.ID)
	if err != nil {
		return
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

func (i *Implement) UpdateOrderItemsStatusServed(ctx context.Context, payload shareModel.OrderItemsStatus) (err error) {
	err = i.IsOrderWithItemsExists(ctx, payload.OrderID, payload.ID)
	if err != nil {
		return
	}

	err = i.repository.UpdateOrderItemsStatusServed(ctx, payload.ID)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to update order items status served", err)
	}

	return
}

func (i *Implement) SearchOrderItems(ctx context.Context, payload domain.SearchOrderItems) (result domain.SearchOrderItemsResult, err error) {
	searchParams := buildSearchOrderItemsParams(payload)

	var (
		searchResult  []*database.SearchOrderItemsRow
		searchErr     error
		totalItems    int64
		totalItemsErr error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		searchResult, searchErr = i.fetchSearchOrder(ctx, searchParams)
	}()

	go func() {
		defer wg.Done()
		totalItems, totalItemsErr = i.fetchSearchOrderTotalItems(ctx, searchParams)
	}()

	wg.Wait()

	if searchErr != nil {
		return domain.SearchOrderItemsResult{}, searchErr
	}

	if totalItemsErr != nil {
		return domain.SearchOrderItemsResult{}, totalItemsErr
	}

	return domain.SearchOrderItemsResult{
		PageNumber: utils.GetPageNumber(payload.PageNumber),
		PageSize:   utils.GetPageSize(payload.PageSize),
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       shareModel.TransformOrderItemsResults(searchResult),
	}, nil
}

func (i *Implement) fetchSearchOrder(ctx context.Context, params database.SearchOrderItemsParams) ([]*database.SearchOrderItemsRow, error) {
	result, err := i.repository.SearchOrderItems(ctx, params)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch order items", err)
	}
	return result, nil
}

func (i *Implement) fetchSearchOrderTotalItems(ctx context.Context, params database.SearchOrderItemsParams) (int64, error) {
	totalParams := database.GetTotalSearchOrderItemsParams{
		ProductName: params.ProductName,
		TableNumber: params.TableNumber,
		StatusCode:  params.StatusCode,
	}
	totalItems, err := i.repository.GetTotalSearchOrderItems(ctx, totalParams)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch total items", err)
	}

	return totalItems, nil
}

func (i *Implement) GetOrderItems(ctx context.Context, orderID int64, search domain.SearchOrderItems) (result domain.SearchOrderItemsResult, err error) {
	searchParams := buildSearchOrderItemsParams(search)

	if orderID <= 0 {
		return domain.SearchOrderItemsResult{}, exceptions.Error(exceptions.CodeBusiness, exceptions.ErrOrderRequired.Error())
	}

	err = i.IsOrderExist(ctx, orderID)
	if err != nil {
		return domain.SearchOrderItemsResult{}, err
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
			Pagesize:   searchParams.PageSize,
			Pagenumber: searchParams.PageNumber,
		})
	}()
	go func() {
		defer wg.Done()
		totalItems, err = i.fetchTotalOrderWithItems(ctx, orderID)
	}()

	wg.Wait()

	if err != nil {
		return domain.SearchOrderItemsResult{}, err
	}

	return domain.SearchOrderItemsResult{
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

func (i *Implement) GetOrderItemsByID(ctx context.Context, orderID, orderItemsID int64, tableNumber int32) (result *shareModel.OrderItems, err error) {

	err = i.IsOrderWithItemsExists(ctx, orderID, orderItemsID)
	if err != nil {
		return nil, err
	}

	items, err := i.repository.GetOrderWithItemsByID(ctx, database.GetOrderWithItemsByIDParams{
		OrderID:      orderID,
		OrderItemsID: orderItemsID,
	})
	if err != nil {

		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderItemsNotFound.Error())
		}

		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to get order items", err)
	}

	return shareModel.TransformOrderItemsByIDResults(items), nil
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
