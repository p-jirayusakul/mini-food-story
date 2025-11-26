package repository

import (
	"context"
	"errors"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"food-story/table-service/internal/domain"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	_errMsgUpdateTableStatusFailed    = "failed to update table status"
	_errMsgFailToFetchTable           = "failed to fetch table"
	_errMsgFailToFetchTotalItemsTable = "failed to fetch total items table"
)

type TableRow interface {
	GetID() int64
	GetTableNumber() int32
	GetStatus() string
	GetStatusEN() string
	GetStatusCode() string
	GetSeats() int32
	GetOrderID() *int64
	GetExpiresAt() pgtype.Timestamptz
	GetExtendTotalMinutes() int32
}

func (i *Implement) IsTableAvailableOrReserved(ctx context.Context, tableID int64) error {
	isAvailableOrReserved, err := i.repository.IsTableAvailableOrReserved(ctx, tableID)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to check table status", err)
	}

	if !isAvailableOrReserved {
		return exceptions.Errorf(exceptions.CodeBusiness, "table not available or reserved", err)
	}

	return nil
}

func (i *Implement) GetTableNumber(ctx context.Context, tableID int64) (int32, error) {
	data, err := i.repository.GetTableNumber(ctx, tableID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, exceptions.ErrorIDNotFound(exceptions.CodeTableNotFound, tableID)
		}
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to get table number", err)
	}

	return data, nil
}

func (i *Implement) ListTableStatus(ctx context.Context) (result []*domain.Status, err error) {
	data, err := i.repository.ListTableStatus(ctx)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to get table status", err)
	}

	if data == nil {
		return nil, exceptions.ErrorDataNotFound()
	}

	result = make([]*domain.Status, len(data))
	for index, v := range data {
		result[index] = &domain.Status{
			ID:     v.ID,
			Code:   v.Code,
			Name:   v.Name,
			NameEn: v.NameEn,
		}
	}

	return result, nil
}

func (i *Implement) SearchTables(ctx context.Context, search domain.SearchTables) (domain.SearchTablesResult, error) {
	searchParams := buildSearchParams(search)

	var (
		searchResult  []*database.SearchTablesRow
		searchErr     error
		totalItems    int64
		totalItemsErr error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		searchResult, searchErr = i.fetchTables(ctx, searchParams)
	}()

	go func() {
		defer wg.Done()
		totalItems, totalItemsErr = i.fetchTotalItems(ctx, searchParams)
	}()

	wg.Wait()

	if searchErr != nil {
		return domain.SearchTablesResult{}, searchErr
	}

	if totalItemsErr != nil {
		return domain.SearchTablesResult{}, totalItemsErr
	}

	return domain.SearchTablesResult{
		PageNumber: utils.GetPageNumber(search.PageNumber),
		PageSize:   utils.GetPageSize(search.PageSize),
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       transformSearchResults(searchResult),
	}, nil
}

func (i *Implement) QuickSearchAvailableTable(ctx context.Context, search domain.SearchTables) (domain.SearchTablesResult, error) {
	searchParams := buildQuickSearchParams(search)

	var (
		searchResult  []*database.QuickSearchTablesRow
		searchErr     error
		totalItems    int64
		totalItemsErr error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		searchResult, searchErr = i.fetchQuickSearchTables(ctx, searchParams)
	}()

	go func() {
		defer wg.Done()
		totalItems, totalItemsErr = i.fetchQuickSearchTablesTotalItems(ctx, search.NumberOfPeople)
	}()

	wg.Wait()

	if searchErr != nil {
		return domain.SearchTablesResult{}, searchErr
	}

	if totalItemsErr != nil {
		return domain.SearchTablesResult{}, totalItemsErr
	}

	return domain.SearchTablesResult{
		PageNumber: utils.GetPageNumber(search.PageNumber),
		PageSize:   utils.GetPageSize(search.PageSize),
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       transformSearchResults(searchResult),
	}, nil
}

func (i *Implement) UpdateTablesStatusAvailable(ctx context.Context, tableID int64) (err error) {
	err = i.repository.UpdateTablesStatusAvailable(ctx, tableID)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, _errMsgUpdateTableStatusFailed, err)
	}
	return nil
}

func (i *Implement) fetchTables(ctx context.Context, params database.SearchTablesParams) ([]*database.SearchTablesRow, error) {
	result, err := i.repository.SearchTables(ctx, params)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, _errMsgFailToFetchTable, err)
	}
	return result, nil
}

func (i *Implement) fetchTotalItems(ctx context.Context, params database.SearchTablesParams) (int64, error) {
	totalParams := database.GetTotalPageSearchTablesParams{
		TableNumber: params.TableNumber,
		Seats:       params.Seats,
		StatusCode:  params.StatusCode,
	}
	totalItems, err := i.repository.GetTotalPageSearchTables(ctx, totalParams)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, _errMsgFailToFetchTotalItemsTable, err)
	}

	return totalItems, nil
}

func (i *Implement) fetchQuickSearchTables(ctx context.Context, params database.QuickSearchTablesParams) ([]*database.QuickSearchTablesRow, error) {
	result, err := i.repository.QuickSearchTables(ctx, params)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, _errMsgFailToFetchTable, err)
	}
	return result, nil
}

func (i *Implement) fetchQuickSearchTablesTotalItems(ctx context.Context, numberOfPeople int32) (int64, error) {
	totalItems, err := i.repository.GetTotalPageQuickSearchTables(ctx, numberOfPeople)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, _errMsgFailToFetchTotalItemsTable, err)
	}

	return totalItems, nil
}

func transformSearchResults[T TableRow](results []T) []*domain.Table {
	data := make([]*domain.Table, len(results))
	for index, row := range results {
		var expiredAt *string
		expiredAtDB, expErr := utils.PgTimestampToThaiISO8601(row.GetExpiresAt())
		if expErr == nil {
			expiredAt = &expiredAtDB
		}
		data[index] = &domain.Table{
			ID:                 row.GetID(),
			TableNumber:        row.GetTableNumber(),
			Status:             row.GetStatus(),
			StatusEn:           row.GetStatusEN(),
			StatusCode:         row.GetStatusCode(),
			Seats:              row.GetSeats(),
			OrderID:            row.GetOrderID(),
			ExpiredAt:          expiredAt,
			ExtendTotalMinutes: row.GetExtendTotalMinutes(),
		}
	}
	return data
}

func buildSearchParams(payload domain.SearchTables) database.SearchTablesParams {
	params := database.SearchTablesParams{
		OrderByType: payload.OrderByType,
		OrderBy:     payload.OrderBy,
		PageSize:    payload.PageSize,
		PageNumber:  payload.PageNumber,
	}

	for _, v := range payload.StatusCode {
		params.StatusCode = append(params.StatusCode, strings.ToUpper(v))
	}

	params.PageSize, params.PageNumber = utils.CalculatePageSizeAndNumber(payload.PageSize, payload.PageNumber)

	if payload.TableNumber != nil {
		params.TableNumber = pgtype.Int4{Int32: *payload.TableNumber, Valid: true}
	}
	if payload.Seats != nil {
		params.Seats = pgtype.Int4{Int32: *payload.Seats, Valid: true}
	}
	return params
}

func buildQuickSearchParams(payload domain.SearchTables) database.QuickSearchTablesParams {
	pageSize, pageNumber := utils.CalculatePageSizeAndNumber(payload.PageSize, payload.PageNumber)
	return database.QuickSearchTablesParams{
		NumberOfPeople: payload.NumberOfPeople,
		OrderByType:    payload.OrderByType,
		OrderBy:        payload.OrderBy,
		PageSize:       pageSize,
		PageNumber:     pageNumber,
	}
}
