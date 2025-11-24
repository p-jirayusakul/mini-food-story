package usecase

import (
	"context"
	"food-story/menu-service/internal/adapter/cache"
	"food-story/menu-service/internal/adapter/repository"
	"food-story/menu-service/internal/domain"
	"food-story/shared/config"

	"github.com/google/uuid"
)

type Usecase interface {
	ListCategory(ctx context.Context) (result []*domain.Category, err error)
	SearchProductByFilters(ctx context.Context, payload domain.SearchProduct) (result domain.SearchProductResult, err error)
	GetProductByID(ctx context.Context, id int64) (result *domain.Product, err error)
	IsSessionValid(sessionID uuid.UUID) error
}

type Implement struct {
	config     config.Config
	repository repository.Implement
	cache      cache.RedisTableCacheInterface
}

func NewUsecase(config config.Config, repository repository.Implement, cache cache.RedisTableCacheInterface) *Implement {
	return &Implement{
		config,
		repository,
		cache,
	}
}

var _ Usecase = (*Implement)(nil)
