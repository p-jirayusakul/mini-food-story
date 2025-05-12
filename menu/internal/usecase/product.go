package usecase

import (
	"context"
	"food-story/menu/internal/domain"
	"food-story/pkg/exceptions"
)

func (i *ProductImplement) ListCategory(ctx context.Context) (result []*domain.Category, customError *exceptions.CustomError) {
	return i.repository.ListCategory(ctx)
}

func (i *ProductImplement) CreateProduct(ctx context.Context, payload domain.Product) (result int64, customError *exceptions.CustomError) {
	return i.repository.CreateProduct(ctx, payload)
}

func (i *ProductImplement) UpdateProduct(ctx context.Context, payload domain.Product) (customError *exceptions.CustomError) {

	customError = i.repository.IsProductExists(ctx, payload.ID)
	if customError != nil {
		return
	}

	return i.repository.UpdateProduct(ctx, payload)
}

func (i *ProductImplement) UpdateProductIsAvailable(ctx context.Context, payload domain.Product) (customError *exceptions.CustomError) {
	customError = i.repository.IsProductExists(ctx, payload.ID)
	if customError != nil {
		return
	}

	return i.repository.UpdateProductAvailability(ctx, payload.ID, payload.IsAvailable)
}

func (i *ProductImplement) SearchProductByFilters(ctx context.Context, payload domain.SearchProduct) (result domain.SearchProductResult, customError *exceptions.CustomError) {
	return i.repository.SearchProduct(ctx, payload)
}
