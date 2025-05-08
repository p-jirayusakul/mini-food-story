package usecase

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/shared/database/sqlc"
	"food-story/shared/snowflakeid"
	"food-story/table/internal/config"
	"food-story/table/internal/domain"
	"github.com/google/uuid"
)

type TableUsecase interface {
	ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError)
	CreateTable(ctx context.Context, payload domain.Table) (result int64, customError *exceptions.CustomError)
	UpdateTable(ctx context.Context, payload domain.Table) (customError *exceptions.CustomError)
	UpdateTableStatus(ctx context.Context, payload domain.TableStatus) (customError *exceptions.CustomError)
	SearchTableByFilters(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError)
	QuickSearchAvailableTable(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError)
	CreateTableSession(ctx context.Context, payload domain.TableSession) (string, *exceptions.CustomError)
	GettableSession(ctx context.Context, sessionID uuid.UUID) (*domain.CurrentTableSession, *exceptions.CustomError)
}

type TableImplement struct {
	config      config.Config
	repository  database.Store
	snowflakeID snowflakeid.SnowflakeInterface
}

func NewUsecase(config config.Config, repository database.Store, snowflakeID snowflakeid.SnowflakeInterface) *TableImplement {
	return &TableImplement{
		config,
		repository,
		snowflakeID,
	}
}

var _ TableUsecase = (*TableImplement)(nil)
