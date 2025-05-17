package usecase

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/table-service/internal/domain"
)

func (i *TableImplement) GetCurrentTime(ctx context.Context) (result domain.TestTime, customError *exceptions.CustomError) {
	return i.repository.GetCurrentTime(ctx)
}
