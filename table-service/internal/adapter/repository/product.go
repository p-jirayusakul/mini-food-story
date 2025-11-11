package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListProductTimeExtension(ctx context.Context) (result []*domain.Product, customError *exceptions.CustomError) {
	data, err := i.repository.ListProductTimeExtension(ctx)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch product time extension: %w", err),
		}
	}

	return transformListProductTimeExtension(data), nil
}

func (i *Implement) GetDurationMinutesByProductID(ctx context.Context, productID int64) (durationMinutes int32, customError *exceptions.CustomError) {
	data, err := i.repository.GetDurationMinutesByProductID(ctx, productID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: errors.New("product id not found"),
			}
		}
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch duration minutes by product id: %w", err),
		}
	}
	return data, nil
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
