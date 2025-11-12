package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"food-story/table-service/internal/domain"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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

func (i *Implement) IsTableAvailableOrReserved(ctx context.Context, tableID int64) *exceptions.CustomError {
	isAvailableOrReserved, err := i.repository.IsTableAvailableOrReserved(ctx, tableID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check table status: %w", err),
		}
	}

	if !isAvailableOrReserved {
		return &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: errors.New("table not available or reserved"),
		}
	}

	return nil
}

func (i *Implement) IsTableExists(ctx context.Context, id int64) *exceptions.CustomError {
	isTableExists, err := i.repository.IsTableExists(ctx, id)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check table exists: %w", err),
		}
	}

	if !isTableExists {
		return &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: errors.New("table not found"),
		}
	}

	return nil
}

func (i *Implement) GetTableNumber(ctx context.Context, tableID int64) (int32, *exceptions.CustomError) {
	data, err := i.repository.GetTableNumber(ctx, tableID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: errors.New("table id not found"),
			}
		}
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get table number: %w", err),
		}
	}

	return data, nil
}

func (i *Implement) ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError) {
	data, err := i.repository.ListTableStatus(ctx)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch table status: %w", err),
		}
	}

	if data == nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: errors.New("table status not found"),
		}
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

func (i *Implement) CreateTable(ctx context.Context, table domain.Table) (result int64, customError *exceptions.CustomError) {
	tableParams := database.CreateTableParams{
		ID:          i.snowflakeID.Generate(),
		TableNumber: table.TableNumber,
		Seats:       table.Seats,
	}

	result, err := i.repository.CreateTable(ctx, tableParams)
	if err != nil {
		return 0, handleUpsertError(err, "tables")
	}

	return result, nil
}

func (i *Implement) UpdateTables(ctx context.Context, table domain.Table) (customError *exceptions.CustomError) {

	tableParams := database.UpdateTablesParams{
		ID:          table.ID,
		TableNumber: table.TableNumber,
		Seats:       table.Seats,
	}

	err := i.repository.UpdateTables(ctx, tableParams)
	if err != nil {
		return handleUpsertError(err, "tables")
	}

	return nil
}

func (i *Implement) UpdateTablesStatus(ctx context.Context, tableStatus domain.TableStatus) (customError *exceptions.CustomError) {

	tableStatusParams := database.UpdateTablesStatusParams{
		ID:       tableStatus.ID,
		StatusID: tableStatus.StatusID,
	}

	err := i.repository.UpdateTablesStatus(ctx, tableStatusParams)
	if err != nil {
		return handleUpsertError(err, "tables")
	}

	return nil
}

func (i *Implement) SearchTables(ctx context.Context, search domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
	searchParams := buildSearchParams(search)

	var (
		searchResult  []*database.SearchTablesRow
		searchErr     *exceptions.CustomError
		totalItems    int64
		totalItemsErr *exceptions.CustomError
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
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       transformSearchResults(searchResult),
	}, nil
}

func (i *Implement) QuickSearchAvailableTable(ctx context.Context, search domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
	searchParams := buildQuickSearchParams(search)

	var (
		searchResult  []*database.QuickSearchTablesRow
		searchErr     *exceptions.CustomError
		totalItems    int64
		totalItemsErr *exceptions.CustomError
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
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       transformSearchResults(searchResult),
	}, nil
}

func (i *Implement) UpdateTablesStatusAvailable(ctx context.Context, tableID int64) (customError *exceptions.CustomError) {
	err := i.repository.UpdateTablesStatusAvailable(ctx, tableID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update table status: %w", err),
		}
	}
	return nil
}

func (i *Implement) fetchTables(ctx context.Context, params database.SearchTablesParams) ([]*database.SearchTablesRow, *exceptions.CustomError) {
	result, err := i.repository.SearchTables(ctx, params)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch products: %w", err),
		}
	}
	return result, nil
}

func (i *Implement) fetchTotalItems(ctx context.Context, params database.SearchTablesParams) (int64, *exceptions.CustomError) {
	totalParams := database.GetTotalPageSearchTablesParams{
		TableNumber: params.TableNumber,
		Seats:       params.Seats,
		StatusCode:  params.StatusCode,
	}
	totalItems, err := i.repository.GetTotalPageSearchTables(ctx, totalParams)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch total items: %w", err),
		}
	}

	return totalItems, nil
}

func (i *Implement) fetchQuickSearchTables(ctx context.Context, params database.QuickSearchTablesParams) ([]*database.QuickSearchTablesRow, *exceptions.CustomError) {
	result, err := i.repository.QuickSearchTables(ctx, params)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch products: %w", err),
		}
	}
	return result, nil
}

func (i *Implement) fetchQuickSearchTablesTotalItems(ctx context.Context, numberOfPeople int32) (int64, *exceptions.CustomError) {
	totalItems, err := i.repository.GetTotalPageQuickSearchTables(ctx, numberOfPeople)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch total items: %w", err),
		}
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

func handleUpsertError(err error, tableName string) *exceptions.CustomError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == exceptions.SqlstateUniqueViolation {
		msg := fmt.Sprintf("%s already exists", utils.IndexToFieldName(pgErr.ConstraintName, tableName))
		return &exceptions.CustomError{
			Status: exceptions.ERRDATACONFLICT,
			Errors: errors.New(msg),
		}
	}

	return &exceptions.CustomError{
		Status: exceptions.ERRREPOSITORY,
		Errors: err,
	}
}
