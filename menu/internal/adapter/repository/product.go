package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/menu/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
)

func (i *ProductRepoImplement) SearchProduct(ctx context.Context, payload domain.SearchProduct) (domain.SearchProductResult, *exceptions.CustomError) {
	searchParams := buildSearchProductParams(payload)

	searchResult, err := i.repository.SearchProducts(ctx, searchParams)
	if err != nil {
		return domain.SearchProductResult{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch product: %w", err),
		}
	}

	totalItemsParam := database.GetTotalPageSearchProductsParams{
		Name:        searchParams.Name,
		IsAvailable: searchParams.IsAvailable,
		CategoryID:  searchParams.CategoryID,
	}

	totalItems, err := i.repository.GetTotalPageSearchProducts(ctx, totalItemsParam)
	if err != nil {
		return domain.SearchProductResult{}, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch product: %w", err),
		}
	}

	data := make([]*domain.Product, len(searchResult))
	for index, row := range searchResult {
		var description, imageURL *string
		var floatPrice float64

		if row.Price.Valid {
			priceRaw, _ := row.Price.Float64Value()
			floatPrice = priceRaw.Float64
		}

		if row.Description.Valid {
			description = &row.Description.String
		}

		if row.ImageUrl.Valid {
			imageURL = &row.ImageUrl.String
		}

		data[index] = &domain.Product{
			ID:             row.ID,
			Name:           row.Name,
			NameEN:         row.NameEn,
			CategoryName:   row.CategoryName,
			CategoryNameEN: row.CategoryNameEN,
			CategoryID:     row.Categories,
			Price:          floatPrice,
			Description:    description,
			IsAvailable:    row.IsAvailable,
			ImageURL:       imageURL,
		}
	}

	return domain.SearchProductResult{
		TotalItems: totalItems,
		TotalPages: int64(math.Ceil(float64(totalItems) / float64(searchParams.PageSize))),
		Data:       data,
	}, nil
}

func (i *ProductRepoImplement) GetProductByID(ctx context.Context, id int64) (result *domain.Product, customError *exceptions.CustomError) {
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
			Errors: fmt.Errorf("failed to get product exists: %w", err),
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

func (i *ProductRepoImplement) UpdateProductAvailability(ctx context.Context, id int64, isAvailable bool) *exceptions.CustomError {
	if err := i.repository.UpdateProductAvailability(ctx, database.UpdateProductAvailabilityParams{
		ID:          id,
		IsAvailable: isAvailable,
	}); err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to update product availability: %w", err),
		}
	}

	return nil
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
