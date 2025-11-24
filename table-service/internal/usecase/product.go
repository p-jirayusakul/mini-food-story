package usecase

import (
	"context"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListProductTimeExtension(ctx context.Context) (result []*domain.Product, err error) {
	return i.repository.ListProductTimeExtension(ctx)
}
