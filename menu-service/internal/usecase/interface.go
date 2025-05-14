package usecase

import (
	"context"
	"food-story/menu-service/internal/adapter/repository"
	"food-story/menu-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/shared/config"
)

type MenuUsecase interface {
	ListCategory(ctx context.Context) (result []*domain.Category, customError *exceptions.CustomError)
	SearchProductByFilters(ctx context.Context, payload domain.SearchProduct) (result domain.SearchProductResult, customError *exceptions.CustomError)
	GetProductByID(ctx context.Context, id int64) (result *domain.Product, customError *exceptions.CustomError)
}

type MenuImplement struct {
	config     config.Config
	repository repository.ProductRepoImplement
}

func NewUsecase(config config.Config, repository repository.ProductRepoImplement) *MenuImplement {
	return &MenuImplement{
		config,
		repository,
	}
}

var _ MenuUsecase = (*MenuImplement)(nil)
