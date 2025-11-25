package repository

import (
	"context"
	"food-story/pkg/exceptions"
)

func (i *Implement) UpdateTablesStatusFoodServed(ctx context.Context, tableID int64) (err error) {
	err = i.repository.UpdateTablesStatusFoodServed(ctx, tableID)
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to update table status", err)
	}

	return nil
}
