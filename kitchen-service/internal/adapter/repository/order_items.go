package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	shareModel "food-story/shared/model"
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
	"sync"
)

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload shareModel.OrderItemsStatus) (customError *exceptions.CustomError) {
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

func (i *Implement) UpdateOrderItemsStatusServed(ctx context.Context, payload shareModel.OrderItemsStatus) (customError *exceptions.CustomError) {
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

	var (
		searchResult  []*database.SearchOrderItemsRow
		searchErr     *exceptions.CustomError
		totalItems    int64
		totalItemsErr *exceptions.CustomError
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
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       shareModel.TransformOrderItemsResults(searchResult),
	}, nil
}

func (i *Implement) fetchSearchOrder(ctx context.Context, params database.SearchOrderItemsParams) ([]*database.SearchOrderItemsRow, *exceptions.CustomError) {
	result, err := i.repository.SearchOrderItems(ctx, params)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch products: %w", err),
		}
	}
	return result, nil
}

func (i *Implement) fetchSearchOrderTotalItems(ctx context.Context, params database.SearchOrderItemsParams) (int64, *exceptions.CustomError) {
	totalParams := database.GetTotalSearchOrderItemsParams{
		ProductName: params.ProductName,
		TableNumber: params.TableNumber,
		StatusCode:  params.StatusCode,
	}
	totalItems, err := i.repository.GetTotalSearchOrderItems(ctx, totalParams)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch total items: %w", err),
		}
	}

	return totalItems, nil
}

func (i *Implement) GetOrderItems(ctx context.Context, orderID int64) (result []*shareModel.OrderItems, customError *exceptions.CustomError) {
	customError = i.IsOrderExist(ctx, orderID)
	if customError != nil {
		return nil, customError
	}

	items, err := i.repository.GetOrderWithItems(ctx, orderID)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get order items: %w", err),
		}
	}

	return shareModel.TransformOrderItemsResults(items), nil
}

func (i *Implement) GetOrderItemsByID(ctx context.Context, orderID, orderItemsID int64, tableNumber int32) (result *shareModel.OrderItems, customError *exceptions.CustomError) {

	customError = i.IsOrderWithItemsExists(ctx, orderID, orderItemsID)
	if customError != nil {
		return nil, customError
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
