package usecase

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database2 "food-story/shared/database/sqlc"
	"food-story/table/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
	"strconv"
	"strings"
	"time"
)

const tableName = "tables"

func (i *TableImplement) ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError) {
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

func (i *TableImplement) CreateTable(ctx context.Context, payload domain.Table) (result int64, customError *exceptions.CustomError) {

	result, err := i.repository.CreateTable(ctx, database2.CreateTableParams{
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

func (i *TableImplement) UpdateTable(ctx context.Context, payload domain.Table) (customError *exceptions.CustomError) {

	customError = i.isTableExists(ctx, payload.ID)
	if customError != nil {
		return customError
	}

	err := i.repository.UpdateTables(ctx, database2.UpdateTablesParams{
		ID:          payload.ID,
		TableNumber: payload.TableNumber,
		Seats:       payload.Seats,
	})
	if err != nil {
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
			Errors: fmt.Errorf("failed to create table: %w", err),
		}
	}

	return nil
}

func (i *TableImplement) UpdateTableStatus(ctx context.Context, payload domain.TableStatus) (customError *exceptions.CustomError) {

	customError = i.isTableExists(ctx, payload.ID)
	if customError != nil {
		return customError
	}

	err := i.repository.UpdateTablesStatus(ctx, database2.UpdateTablesStatusParams{
		ID:       payload.ID,
		StatusID: payload.StatusID,
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update table status: %w", err),
		}
	}

	return nil
}

func (i *TableImplement) SearchTableByFilters(ctx context.Context, payload domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
	searchParams := buildSearchTablesParams(payload)

	searchResult, err := i.repository.SearchTables(ctx, searchParams)
	if err != nil {
		return domain.SearchTablesResult{}, buildRepositoryError(err, "failed to fetch table status")
	}

	totalItemsParam := database2.GetTotalPageSearchTablesParams{
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

func (i *TableImplement) QuickSearchAvailableTable(ctx context.Context, payload domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
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

func (i *TableImplement) CreateTableSession(ctx context.Context, payload domain.TableSession) (string, *exceptions.CustomError) {

	isAvailableOrReserved, err := i.repository.IsTableAvailableOrReserved(ctx, payload.TableID)
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check table status: %w", err),
		}
	}

	if !isAvailableOrReserved {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: errors.New("table not available or reserved"),
		}
	}

	// เพิ่ม 1 ชั่วโมง
	expiry := time.Now().Add(1 * time.Hour)
	sessionID, err := i.repository.CreateTableSession(ctx, database2.CreateTableSessionParams{
		ID:             i.snowflakeID.Generate(),
		TableID:        payload.TableID,
		NumberOfPeople: payload.NumberOfPeople,
		ExpireAt:       pgtype.Timestamp{Time: expiry, Valid: true},
	})
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to create table session: %w", err),
		}
	}

	err = i.repository.UpdateTablesStatusOccupied(ctx, payload.TableID)
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update table status: %w", err),
		}
	}

	sessionIDEncrypt, err := utils.EncryptSession(utils.SessionData{
		SessionID: sessionID.String(),
		Expiry:    expiry,
	}, []byte(i.config.SecretKey))
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRUNKNOWN,
			Errors: fmt.Errorf("failed to create table session: %w", err),
		}
	}

	url := "http://localhost:3000?s=" + sessionIDEncrypt

	return url, nil
}

func (i *TableImplement) GettableSession(ctx context.Context, sessionID uuid.UUID) (*domain.CurrentTableSession, *exceptions.CustomError) {
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

func (i *TableImplement) isTableExists(ctx context.Context, id int64) *exceptions.CustomError {
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

func buildSearchTablesParams(payload domain.SearchTables) database2.SearchTablesParams {
	params := database2.SearchTablesParams{
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

func buildQuickSearchParams(payload domain.SearchTables) database2.QuickSearchTablesParams {
	pageSize, pageNumber := utils.CalculatePageSizeAndNumber(payload.PageSize, payload.PageNumber)
	return database2.QuickSearchTablesParams{
		NumberOfPeople: payload.NumberOfPeople,
		OrderByType:    payload.OrderByType,
		OrderBy:        payload.OrderBy,
		PageSize:       pageSize,
		PageNumber:     pageNumber,
	}
}
