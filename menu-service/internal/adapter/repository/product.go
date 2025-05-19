package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/menu-service/internal/domain"
	"food-story/pkg/common"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
)

func (i *Implement) SearchProduct(ctx context.Context, payload domain.SearchProduct) (domain.SearchProductResult, *exceptions.CustomError) {
	searchParams := buildSearchProductParams(payload)

	if ctx.Err() != nil {
		return domain.SearchProductResult{}, &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: exceptions.ErrCtxCanceledOrTimeout,
		}
	}

	searchResult, productFetchErr := i.fetchProducts(ctx, searchParams)
	if productFetchErr != nil {
		return domain.SearchProductResult{}, productFetchErr
	}

	totalItems, totalItemsFetchErr := i.fetchTotalItems(ctx, searchParams)
	if totalItemsFetchErr != nil {
		return domain.SearchProductResult{}, totalItemsFetchErr
	}

	return domain.SearchProductResult{
		TotalItems: totalItems,
		TotalPages: calculateTotalPages(totalItems, searchParams.PageSize),
		Data:       transformSearchResults(searchResult),
	}, nil
}

func (i *Implement) GetProductByID(ctx context.Context, id int64) (*domain.Product, *exceptions.CustomError) {

	if ctx.Err() != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRBUSSINESS,
			Errors: exceptions.ErrCtxCanceledOrTimeout,
		}
	}

	data, err := i.repository.GetProductAvailableByID(ctx, id)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: errors.New("product not found"),
			}
		}
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: errorGetProductFailed(id, err),
		}
	}

	if data == nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: errorGetProductFailed(id, err),
		}
	}

	return &domain.Product{
		ID:             data.ID,
		Name:           data.Name,
		NameEN:         data.NameEn,
		CategoryName:   data.CategoryName,
		CategoryNameEN: data.CategoryNameEN,
		CategoryID:     data.Categories,
		Price:          utils.PgNumericToFloat64(data.Price),
		Description:    utils.PgTextToStringPtr(data.Description),
		IsAvailable:    data.IsAvailable,
		ImageURL:       utils.PgTextToStringPtr(data.ImageUrl),
	}, nil
}

func (i *Implement) fetchProducts(ctx context.Context, params database.SearchProductsParams) ([]*database.SearchProductsRow, *exceptions.CustomError) {
	result, err := i.repository.SearchProducts(ctx, params)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch products: %w", err),
		}
	}
	return result, nil
}

func (i *Implement) fetchTotalItems(ctx context.Context, params database.SearchProductsParams) (int64, *exceptions.CustomError) {
	totalParams := database.GetTotalPageSearchProductsParams{
		Name:        params.Name,
		IsAvailable: params.IsAvailable,
		CategoryID:  params.CategoryID,
	}
	totalItems, err := i.repository.GetTotalPageSearchProducts(ctx, totalParams)
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch total items: %w", err),
		}
	}
	return totalItems, nil
}

func buildSearchProductParams(payload domain.SearchProduct) database.SearchProductsParams {
	params := database.SearchProductsParams{
		Name:        pgtype.Text{String: payload.Name, Valid: payload.Name != ""},
		IsAvailable: pgtype.Bool{Bool: payload.IsAvailable, Valid: true},
		CategoryID:  payload.CategoryID,
		OrderByType: payload.OrderByType,
		OrderBy:     payload.OrderBy,
		PageSize:    payload.PageSize,
		PageNumber:  payload.PageNumber,
	}

	params.PageSize, params.PageNumber = utils.CalculatePageSizeAndNumber(payload.PageSize, payload.PageNumber)

	return params
}

func transformSearchResults(results []*database.SearchProductsRow) []*domain.Product {
	data := make([]*domain.Product, len(results))
	for index, row := range results {

		if row == nil {
			continue
		}

		data[index] = &domain.Product{
			ID:             row.ID,
			Name:           row.Name,
			NameEN:         row.NameEn,
			CategoryName:   row.CategoryName,
			CategoryNameEN: row.CategoryNameEN,
			CategoryID:     row.Categories,
			Price:          utils.PgNumericToFloat64(row.Price),
			Description:    utils.PgTextToStringPtr(row.Description),
			IsAvailable:    row.IsAvailable,
			ImageURL:       utils.PgTextToStringPtr(row.ImageUrl),
		}
	}
	return data
}

func calculateTotalPages(totalItems int64, pageSize int64) int64 {
	if pageSize <= 0 {
		pageSize = common.DefaultPageSize
	}

	return int64(math.Ceil(float64(totalItems) / float64(pageSize)))
}

func errorGetProductFailed(id int64, err error) error {
	return fmt.Errorf("get product failed, id: %d, error: %w", id, err)
}
