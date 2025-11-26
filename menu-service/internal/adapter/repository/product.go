package repository

import (
	"context"
	"errors"
	"food-story/menu-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

func (i *Implement) SearchProduct(ctx context.Context, payload domain.SearchProduct) (domain.SearchProductResult, error) {
	searchParams := buildSearchParams(payload)

	var (
		searchResult  []*database.SearchProductsRow
		searchErr     error
		totalItems    int64
		totalItemsErr error
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
		PageNumber: utils.GetPageNumber(payload.PageNumber),
		PageSize:   utils.GetPageSize(payload.PageSize),
		TotalItems: totalItems,
		TotalPages: utils.CalculateTotalPages(totalItems, searchParams.PageSize),
		Data:       transformSearchResults(searchResult),
	}, nil
}

func (i *Implement) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {

	const _errorGetProductFailed = "failed to get product by id"

	data, err := i.repository.GetProductAvailableByID(ctx, id)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return nil, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrProductNotFound.Error())
		}
		return nil, exceptions.Errorf(exceptions.CodeRepository, _errorGetProductFailed, err)
	}

	if data == nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, _errorGetProductFailed, err)
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

func (i *Implement) fetchProducts(ctx context.Context, params database.SearchProductsParams) ([]*database.SearchProductsRow, error) {
	result, err := i.repository.SearchProducts(ctx, params)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch product", err)
	}
	return result, nil
}

func (i *Implement) fetchTotalItems(ctx context.Context, params database.SearchProductsParams) (int64, error) {
	totalParams := database.GetTotalPageSearchProductsParams{
		Name:        params.Name,
		IsAvailable: params.IsAvailable,
		CategoryID:  params.CategoryID,
	}
	totalItems, err := i.repository.GetTotalPageSearchProducts(ctx, totalParams)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch total items", err)
	}
	return totalItems, nil
}

func (i *Implement) ListProductTimeExtension(ctx context.Context) (result []*domain.Product, err error) {
	data, err := i.repository.ListProductTimeExtension(ctx)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch product time extension", err)
	}

	return transformListProductTimeExtension(data), nil
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

func transformListProductTimeExtension(results []*database.ListProductTimeExtensionRow) []*domain.Product {
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
