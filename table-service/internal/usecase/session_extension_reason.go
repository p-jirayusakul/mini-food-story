package usecase

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListSessionExtensionReason(ctx context.Context) (result []*domain.ListSessionExtensionReason, customError *exceptions.CustomError) {
	return i.repository.ListSessionExtensionReason(ctx)
}
