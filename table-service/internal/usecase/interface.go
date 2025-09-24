package usecase

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/shared/config"
	"food-story/table-service/internal/adapter/cache"
	"food-story/table-service/internal/adapter/repository"
	"food-story/table-service/internal/domain"
)

type UseCase interface {
	ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError)
	CreateTable(ctx context.Context, payload domain.Table) (result int64, customError *exceptions.CustomError)
	UpdateTable(ctx context.Context, payload domain.Table) (customError *exceptions.CustomError)
	UpdateTableStatus(ctx context.Context, payload domain.TableStatus) (customError *exceptions.CustomError)
	SearchTableByFilters(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError)
	QuickSearchAvailableTable(ctx context.Context, payload domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError)
	CreateTableSession(ctx context.Context, payload domain.TableSession) (result string, customError *exceptions.CustomError)
	UpdateTableStatusAvailable(ctx context.Context, tableID int64) (customError *exceptions.CustomError)
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

var _ UseCase = (*Implement)(nil)
