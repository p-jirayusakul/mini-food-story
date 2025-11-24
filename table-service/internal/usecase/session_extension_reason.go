package usecase

import (
	"context"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListSessionExtensionReason(ctx context.Context) (result []*domain.ListSessionExtensionReason, err error) {
	return i.repository.ListSessionExtensionReason(ctx)
}
