package usecase

import (
	"context"
	database "food-story/internal/shared/database/sqlc"
	"food-story/internal/shared/snowflakeid"
	"food-story/internal/table/config"
	"food-story/internal/table/domain"
	"food-story/pkg/exceptions"
)

type Usecase interface {
	ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError)
	CreateTable(ctx context.Context, payload domain.Table) (result int64, customError *exceptions.CustomError)
	UpdateTable(ctx context.Context, payload domain.Table) (customError *exceptions.CustomError)
	UpdateTableStatus(ctx context.Context, payload domain.TableStatus) (customError *exceptions.CustomError)
	SearchTableByFilters(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError)
	QuickSearchAvailableTable(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError)
	CreateTableSession(ctx context.Context, payload domain.TableSession) (string, *exceptions.CustomError)
}

type Implement struct {
	config      config.Config
	repository  database.Store
	snowflakeID snowflakeid.SnowflakeInterface
}

func NewUsecase(config config.Config, repository database.Store, snowflakeID snowflakeid.SnowflakeInterface) *Implement {
	return &Implement{
		config,
		repository,
		snowflakeID,
	}
}

var _ Usecase = (*Implement)(nil)
