package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"food-story/table-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
	"strconv"
	"strings"
	"time"
)

func (i *Implement) ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError) {
	data, err := i.repository.ListTableStatus(ctx)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch table status: %w", err),
		}
	}

	if data == nil {
		return nil, nil
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
		return 0, handleUniqueConstraintError(err, "tables")
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
		return handleUniqueConstraintError(err, "tables")
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
		return handleUniqueConstraintError(err, "tables")
	}

	return nil
}

func (i *Implement) SearchTables(ctx context.Context, search domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
	searchParams := buildSearchParams(search)

	if ctx.Err() != nil {
		return domain.SearchTablesResult{}, &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: exceptions.ErrCtxCanceledOrTimeout,
		}
	}

	searchResult, tablesFetchErr := i.fetchTables(ctx, searchParams)
	if tablesFetchErr != nil {
		return domain.SearchTablesResult{}, tablesFetchErr
	}

	totalItems, totalItemsFetchErr := i.fetchTotalItems(ctx, searchParams)
	if totalItemsFetchErr != nil {
		return domain.SearchTablesResult{}, totalItemsFetchErr
	}

	return domain.SearchTablesResult{
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       transformSearchResults(searchResult),
	}, nil
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

func transformSearchResults(results []*database.SearchTablesRow) []*domain.Table {
	data := make([]*domain.Table, len(results))
	for index, row := range results {

		if row == nil {
			continue
		}

		data[index] = &domain.Table{
			ID:          row.ID,
			TableNumber: row.TableNumber,
			Status:      row.Status,
			StatusEn:    row.StatusEN,
			Seats:       row.Seats,
		}
	}
	return data
}

func (i *Implement) QuickSearchTables(ctx context.Context, payload domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
	searchParams := buildQuickSearchParams(payload)

	searchResult, err := i.repository.QuickSearchTables(ctx, searchParams)
	if err != nil {
		return domain.SearchTablesResult{}, buildRepositoryError(err, "failed to fetch table status")
	}

	totalItems, err := i.repository.GetTotalPageQuickSearchTables(ctx, payload.NumberOfPeople)
	if err != nil {
		return domain.SearchTablesResult{}, buildRepositoryError(err, "failed to fetch total page count")
	}

	data := make([]*domain.Table, len(searchResult))
	for index, row := range searchResult {
		data[index] = &domain.Table{
			ID:          row.ID,
			TableNumber: row.TableNumber,
			Status:      row.Status,
			StatusEn:    row.StatusEN,
			Seats:       row.Seats,
		}
	}

	return domain.SearchTablesResult{
		TotalItems: totalItems,
		TotalPages: int64(math.Ceil(float64(totalItems) / float64(searchParams.PageSize))),
		Data:       data,
	}, nil
}

func (i *Implement) CreateTableSession(ctx context.Context, payload domain.TableSession, sessionID uuid.UUID, expiry time.Time) *exceptions.CustomError {

	var sessionByte [16]byte = sessionID
	err := i.repository.TXCreateTableSession(ctx, database.CreateTableSessionParams{
		ID:             i.snowflakeID.Generate(),
		TableID:        payload.TableID,
		NumberOfPeople: payload.NumberOfPeople,
		SessionID:      pgtype.UUID{Bytes: sessionByte, Valid: true},
		ExpireAt:       pgtype.Timestamptz{Time: expiry, Valid: true},
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to create table session: %w", err),
		}
	}

	return nil
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

func (i *Implement) GettableSession(ctx context.Context, sessionID uuid.UUID) (*domain.CurrentTableSession, *exceptions.CustomError) {
	var byteArray [16]byte = sessionID
	id := pgtype.UUID{
		Bytes: byteArray,
		Valid: true,
	}

	isExists, err := i.repository.IsTableSessionExists(ctx, id)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check table session exists: %w", err),
		}
	}

	if !isExists {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: errors.New("table session not found"),
		}
	}

	data, err := i.repository.GetTableSession(ctx, id)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get table session: %w", err),
		}
	}

	var result domain.CurrentTableSession
	result = domain.CurrentTableSession{
		SessionID:   sessionID,
		TableID:     data.TableID,
		TableNumber: data.TableNumber,
		Status:      string(data.Status.TableSessionStatus),
		StartedAt:   data.StartedAt.Time,
	}

	if data.OrderID.Valid {
		orderID := strconv.FormatInt(data.OrderID.Int64, 10)
		result.OrderID = &orderID
	}

	return &result, nil
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
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to get table number: %w", err),
		}
	}

	return data, nil
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

func buildRepositoryError(err error, msg string) *exceptions.CustomError {
	return &exceptions.CustomError{
		Status: exceptions.ERRREPOSITORY,
		Errors: fmt.Errorf("%s: %w", msg, err),
	}
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

func handleUniqueConstraintError(err error, tableName string) *exceptions.CustomError {
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
