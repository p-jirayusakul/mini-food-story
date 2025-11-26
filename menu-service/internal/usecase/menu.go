package usecase

import (
	"context"
	"food-story/menu-service/internal/domain"
)

func (i *Implement) ListCategory(ctx context.Context) (result []*domain.Category, err error) {
	return i.repository.ListCategory(ctx)
}

func (i *Implement) SearchProductByFilters(ctx context.Context, payload domain.SearchProduct) (result domain.SearchProductResult, err error) {
	return i.repository.SearchProduct(ctx, payload)
}

func (i *Implement) GetProductByID(ctx context.Context, id int64) (result *domain.Product, err error) {
	return i.repository.GetProductByID(ctx, id)
}

func (i *Implement) ListProductTimeExtension(ctx context.Context) (result []*domain.Product, err error) {
	return i.repository.ListProductTimeExtension(ctx)
}
