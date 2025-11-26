package usecase

import (
	"context"
	"food-story/shared/config"
	"food-story/table-service/internal/adapter/cache"
	"food-story/table-service/internal/adapter/repository"
	"food-story/table-service/internal/domain"
)

type UseCase interface {
	ListTableStatus(ctx context.Context) (result []*domain.Status, err error)
	UpdateTableStatus(ctx context.Context, payload domain.TableStatus) (err error)
	SearchTableByFilters(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, err error)
	QuickSearchAvailableTable(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, err error)
	CreateTableSession(ctx context.Context, payload domain.TableSession) (result string, err error)
	UpdateTableStatusAvailable(ctx context.Context, tableID int64) (err error)
	ListSessionExtensionReason(ctx context.Context) (result []*domain.ListSessionExtensionReason, err error)
	SessionExtension(ctx context.Context, payload domain.SessionExtension) error
}

type Implement struct {
	config     config.Config
	repository repository.Implement
	cache      cache.RedisTableCacheInterface
}

func NewUseCase(config config.Config, repository repository.Implement, cache cache.RedisTableCacheInterface) *Implement {
	return &Implement{
		config,
		repository,
		cache,
	}
}

var _ UseCase = (*Implement)(nil)
