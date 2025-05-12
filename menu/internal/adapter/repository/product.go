package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/menu/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
	"math/big"
)

func (i *ProductRepoImplement) IsProductExists(ctx context.Context, id int64) *exceptions.CustomError {
	isProductExists, err := i.repository.IsProductExists(ctx, id)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check product exists: %w", err),
		}
	}

	if !isProductExists {
		return &exceptions.CustomError{
			Status: exceptions.ERRNOTFOUND,
			Errors: errors.New("product not found"),
		}
	}

	return nil
}

func (i *ProductRepoImplement) CreateProduct(ctx context.Context, payload domain.Product) (int64, *exceptions.CustomError) {
	var description pgtype.Text
	var imageURL pgtype.Text
	var price pgtype.Numeric

	if payload.Description != nil {
		description.String = *payload.Description
		description.Valid = true
	}

	if payload.ImageURL != nil {
		imageURL.String = *payload.ImageURL
		imageURL.Valid = true
	}

	price = pgtype.Numeric{Int: big.NewInt(utils.ConvertFloatToIntExp(payload.Price)), Exp: -2, Valid: true}

	params := database.CreateProductParams{
		ID:          i.snowflakeID.Generate(),
		Name:        payload.Name,
		NameEn:      payload.NameEN,
		Categories:  payload.CategoryID,
		Description: description,
		Price:       price,
		IsAvailable: payload.IsAvailable,
		ImageUrl:    imageURL,
	}

	id, err := i.repository.CreateProduct(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == exceptions.SqlstateUniqueViolation {
			msg := fmt.Sprintf("%s already exists", utils.IndexToFieldName(pgErr.ConstraintName, "products"))
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

	return id, nil
}

func (i *ProductRepoImplement) UpdateProduct(ctx context.Context, payload domain.Product) *exceptions.CustomError {
	var description pgtype.Text
	var imageURL pgtype.Text
	var price pgtype.Numeric

	if payload.Description != nil {
		description.String = *payload.Description
		description.Valid = true
	}

	if payload.ImageURL != nil {
		imageURL.String = *payload.ImageURL
		imageURL.Valid = true
	}

	price = pgtype.Numeric{Int: big.NewInt(utils.ConvertFloatToIntExp(payload.Price)), Exp: -2, Valid: true}

	params := database.UpdateProductParams{
		ID:          payload.ID,
		Name:        payload.Name,
		NameEn:      payload.NameEN,
		Categories:  payload.CategoryID,
		Description: description,
		Price:       price,
		IsAvailable: payload.IsAvailable,
		ImageUrl:    imageURL,
	}

	if err := i.repository.UpdateProduct(ctx, params); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == exceptions.SqlstateUniqueViolation {
			msg := fmt.Sprintf("%s already exists", utils.IndexToFieldName(pgErr.ConstraintName, "products"))
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
			ID:          row.ID,
			Name:        row.Name,
			NameEN:      row.NameEn,
			CategoryID:  row.Categories,
			Price:       floatPrice,
			Description: description,
			IsAvailable: row.IsAvailable,
			ImageURL:    imageURL,
		}
	}

	return domain.SearchProductResult{
		TotalItems: totalItems,
		TotalPages: int64(math.Ceil(float64(totalItems) / float64(searchParams.PageSize))),
		Data:       data,
	}, nil
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
