package usecase

import (
	"context"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListTableStatus(ctx context.Context) (result []*domain.Status, err error) {
	return i.repository.ListTableStatus(ctx)
}

func (i *Implement) UpdateTableStatus(ctx context.Context, tableStatus domain.TableStatus) (err error) {

	if err := i.repository.IsTableExists(ctx, tableStatus.ID); err != nil {
		return err
	}

	return i.repository.UpdateTablesStatus(ctx, tableStatus)
}

func (i *Implement) UpdateTableStatusAvailable(ctx context.Context, tableID int64) (err error) {
	return i.repository.UpdateTablesStatusAvailable(ctx, tableID)
}

func (i *Implement) SearchTableByFilters(ctx context.Context, search domain.SearchTables) (result domain.SearchTablesResult, err error) {
	return i.repository.SearchTables(ctx, search)
}

func (i *Implement) QuickSearchAvailableTable(ctx context.Context, search domain.SearchTables) (domain.SearchTablesResult, error) {
	return i.repository.QuickSearchAvailableTable(ctx, search)
}
