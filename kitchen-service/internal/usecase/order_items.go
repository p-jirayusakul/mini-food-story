package usecase

import (
	"context"
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
)

func (i *Implement) UpdateOrderItemsStatus(ctx context.Context, payload domain.OrderItemsStatus) (customError *exceptions.CustomError) {
	return i.repository.UpdateOrderItemsStatus(ctx, payload)
}
