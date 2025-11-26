package usecase

import (
	"context"
	"errors"
	"food-story/pkg/exceptions"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListTableStatus(ctx context.Context) (result []*domain.Status, err error) {
	return i.getTableStatusFromCache(ctx)
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

func (i *Implement) getTableStatusFromCache(ctx context.Context) ([]*domain.Status, error) {
	tableStatus, err := i.cache.GetCachedTableStatus()
	if err != nil {
		if errors.Is(err, exceptions.ErrRedisKeyNotFoundException) {
			tableStatusDB, getTableNumberErr := i.repository.ListTableStatus(ctx)
			if getTableNumberErr != nil {
				return nil, getTableNumberErr
			}

			setCacheErr := i.cache.SetCachedTableStatus(tableStatusDB)
			if setCacheErr != nil {
				return nil, setCacheErr
			}

			return tableStatusDB, nil
		}
		return nil, err
	}

	return tableStatus, nil
}
