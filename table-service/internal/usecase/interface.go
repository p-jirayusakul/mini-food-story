package usecase

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/shared/config"
	"food-story/table-service/internal/adapter/cache"
	"food-story/table-service/internal/adapter/repository"
	"food-story/table-service/internal/domain"
)

type TableUsecase interface {
	ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError)
	CreateTable(ctx context.Context, payload domain.Table) (result int64, customError *exceptions.CustomError)
	UpdateTable(ctx context.Context, payload domain.Table) (customError *exceptions.CustomError)
	UpdateTableStatus(ctx context.Context, payload domain.TableStatus) (customError *exceptions.CustomError)
	SearchTableByFilters(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError)
	QuickSearchAvailableTable(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError)
	CreateTableSession(ctx context.Context, payload domain.TableSession) (string, *exceptions.CustomError)
	GetCurrentSession(sessionIDEncrypt string) (*domain.CurrentTableSession, *exceptions.CustomError)

	GetCurrentTime(ctx context.Context) (result domain.TestTime, customError *exceptions.CustomError)
}

type TableImplement struct {
	config     config.Config
	repository repository.TableRepoImplement
	cache      cache.RedisTableCacheInterface
}

func NewUsecase(config config.Config, repository repository.TableRepoImplement, cache cache.RedisTableCacheInterface) *TableImplement {
	return &TableImplement{
		config,
		repository,
		cache,
	}
}

var _ TableUsecase = (*TableImplement)(nil)
