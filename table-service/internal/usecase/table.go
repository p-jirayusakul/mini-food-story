package usecase

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListTableStatus(ctx context.Context) (result []*domain.Status, customError *exceptions.CustomError) {
	return i.repository.ListTableStatus(ctx)
}

func (i *Implement) CreateTable(ctx context.Context, payload domain.Table) (result int64, customError *exceptions.CustomError) {
	return i.repository.CreateTable(ctx, payload)
}

func (i *Implement) UpdateTable(ctx context.Context, payload domain.Table) (customError *exceptions.CustomError) {

	customError = i.repository.IsTableExists(ctx, payload.ID)
	if customError != nil {
		return
	}

	return i.repository.UpdateTables(ctx, payload)
}

func (i *Implement) UpdateTableStatus(ctx context.Context, payload domain.TableStatus) (customError *exceptions.CustomError) {

	customError = i.repository.IsTableExists(ctx, payload.ID)
	if customError != nil {
		return
	}

	customError = i.repository.UpdateTablesStatus(ctx, payload)
	if customError != nil {
		return
	}

	return nil
}

func (i *Implement) SearchTableByFilters(ctx context.Context, payload domain.SearchTables) (result domain.SearchTablesResult, customError *exceptions.CustomError) {
	return i.repository.SearchTables(ctx, payload)
}

func (i *Implement) QuickSearchAvailableTable(ctx context.Context, payload domain.SearchTables) (domain.SearchTablesResult, *exceptions.CustomError) {
	return i.repository.QuickSearchTables(ctx, payload)
}
