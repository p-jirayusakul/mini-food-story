package usecase

import (
	"context"
	"errors"
	"fmt"
	database "food-story/internal/shared/database/sqlc"
	"food-story/internal/table/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
	"strings"
)

const tableName = "tables"

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

func (i *Implement) CreateTable(ctx context.Context, payload domain.CreateTableParam) (result int64, customError *exceptions.CustomError) {

	result, err := i.repository.CreateTable(ctx, database.CreateTableParams{
		ID:          i.snowflakeID.Generate(),
		TableNumber: payload.TableNumber,
		Seats:       payload.Seats,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == exceptions.SqlstateUniqueViolation {
			msg := fmt.Sprintf("%s already exists", utils.IndexToFieldName(pgErr.ConstraintName, tableName))
			return 0, &exceptions.CustomError{
				Status: exceptions.ERRDATACONFLICT,
				Errors: errors.New(msg),
			}
		}

		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to create table: %w", err),
		}
	}

	return result, nil
}

func (i *Implement) SearchTableByFilters(ctx context.Context, payload domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
	searchParams := buildSearchTablesParams(payload)

	searchResult, err := i.repository.SearchTables(ctx, searchParams)
	if err != nil {
		return domain.SearchTablesResult{}, buildRepositoryError(err, "failed to fetch table status")
	}

	totalItemsParam := database.GetTotalPageSearchTablesParams{
		TableNumber: searchParams.TableNumber,
		Seats:       searchParams.Seats,
		StatusCode:  searchParams.StatusCode,
	}

	totalItems, err := i.repository.GetTotalPageSearchTables(ctx, totalItemsParam)
	if err != nil {
		return domain.SearchTablesResult{}, buildRepositoryError(err, "failed to fetch table status")
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
		TotalPages: int64(math.Ceil(float64(totalItems) / float64(payload.PageSize))),
		Data:       data,
	}, nil
}

func (i *Implement) QuickSearchAvailableTable(ctx context.Context, payload domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
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
		TotalPages: int64(math.Ceil(float64(totalItems) / float64(payload.PageSize))),
		Data:       data,
	}, nil
}

func buildSearchTablesParams(payload domain.SearchTables) database.SearchTablesParams {
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
