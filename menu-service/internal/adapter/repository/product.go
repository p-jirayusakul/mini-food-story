package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/menu-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

func (i *Implement) SearchProduct(ctx context.Context, payload domain.SearchProduct) (domain.SearchProductResult, *exceptions.CustomError) {
	searchParams := buildSearchParams(payload)

	var (
		searchResult  []*database.SearchProductsRow
		searchErr     *exceptions.CustomError
		totalItems    int64
		totalItemsErr *exceptions.CustomError
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		searchResult, searchErr = i.fetchProducts(ctx, searchParams)
	}()

	go func() {
		defer wg.Done()
		totalItems, totalItemsErr = i.fetchTotalItems(ctx, searchParams)
	}()

	wg.Wait()

	if searchErr != nil {
		return domain.SearchProductResult{}, searchErr
	}

	if totalItemsErr != nil {
		return domain.SearchProductResult{}, totalItemsErr
	}

	return domain.SearchProductResult{
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       transformSearchResults(searchResult),
	}, nil
}

func (i *Implement) GetProductByID(ctx context.Context, id int64) (*domain.Product, *exceptions.CustomError) {

	data, err := i.repository.GetProductAvailableByID(ctx, id)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrProductNotFound,
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

func buildSearchParams(payload domain.SearchProduct) database.SearchProductsParams {
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

func errorGetProductFailed(id int64, err error) error {
	return fmt.Errorf("get product failed, id: %d, error: %w", id, err)
}
