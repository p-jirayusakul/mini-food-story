package usecase

import (
	"context"
	"food-story/menu/internal/adapter/repository"
	"food-story/menu/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/shared/config"
)

type ProductUsecase interface {
	ListCategory(ctx context.Context) (result []*domain.Category, customError *exceptions.CustomError)
	SearchProductByFilters(ctx context.Context, payload domain.SearchProduct) (result domain.SearchProductResult, customError *exceptions.CustomError)
	GetProductByID(ctx context.Context, id int64) (result *domain.Product, customError *exceptions.CustomError)
}

type ProductImplement struct {
	config     config.Config
	repository repository.ProductRepoImplement
}

func NewUsecase(config config.Config, repository repository.ProductRepoImplement) *ProductImplement {
	return &ProductImplement{
		config,
		repository,
	}
}

var _ ProductUsecase = (*ProductImplement)(nil)
