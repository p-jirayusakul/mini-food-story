package usecase

import (
	"context"
	"food-story/menu/internal/domain"
	"food-story/pkg/exceptions"
)

func (i *ProductImplement) ListCategory(ctx context.Context) (result []*domain.Category, customError *exceptions.CustomError) {
	return i.repository.ListCategory(ctx)
}

func (i *ProductImplement) SearchProductByFilters(ctx context.Context, payload domain.SearchProduct) (result domain.SearchProductResult, customError *exceptions.CustomError) {
	return i.repository.SearchProduct(ctx, payload)
}

func (i *ProductImplement) GetProductByID(ctx context.Context, id int64) (result *domain.Product, customError *exceptions.CustomError) {
	return i.repository.GetProductByID(ctx, id)
}
