package usecase

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError) {
	return i.repository.ListTableStatus(ctx)
}

func (i *Implement) CreateTable(ctx context.Context, table domain.Table) (result int64, customError *exceptions.CustomError) {
	return i.repository.CreateTable(ctx, table)
}

func (i *Implement) UpdateTable(ctx context.Context, table domain.Table) (customError *exceptions.CustomError) {

	customError = i.repository.IsTableExists(ctx, table.ID)
	if customError != nil {
		return customError
	}

	return i.repository.UpdateTables(ctx, table)
}

func (i *Implement) UpdateTableStatus(ctx context.Context, tableStatus domain.TableStatus) (customError *exceptions.CustomError) {

	customError = i.repository.IsTableExists(ctx, tableStatus.ID)
	if customError != nil {
		return customError
	}

	return i.repository.UpdateTablesStatus(ctx, tableStatus)
}

func (i *Implement) SearchTableByFilters(ctx context.Context, search domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError) {
	return i.repository.SearchTables(ctx, search)
}

func (i *Implement) QuickSearchAvailableTable(ctx context.Context, search domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
	return i.repository.QuickSearchAvailableTable(ctx, search)
}
