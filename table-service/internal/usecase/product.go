package usecase

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListProductTimeExtension(ctx context.Context) (result []*domain.Product, customError *exceptions.CustomError) {
	return i.repository.ListProductTimeExtension(ctx)
}
