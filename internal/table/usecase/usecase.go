package usecase

import (
	"context"
	database "food-story/internal/shared/database/sqlc"
	"food-story/internal/shared/snowflakeid"
	"food-story/internal/table/domain"
	"food-story/pkg/exceptions"
)

type Usecase interface {
	ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError)
	CreateTable(ctx context.Context, payload domain.CreateTableParam) (result int64, customError *exceptions.CustomError)
	SearchTableByFilters(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError)
	QuickSearchAvailableTable(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError)
}

type Implement struct {
	repository  database.Store
	snowflakeID snowflakeid.SnowflakeInterface
}

func NewUsecase(repository database.Store, snowflakeID snowflakeid.SnowflakeInterface) *Implement {
	return &Implement{
		repository,
		snowflakeID,
	}
}

var _ Usecase = (*Implement)(nil)
