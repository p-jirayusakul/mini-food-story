package usecase

import (
	"context"
	"food-story/payment-service/internal/domain"
)

func (i *PaymentImplement) ListPaymentMethods(ctx context.Context) (result []*domain.PaymentMethod, err error) {
	return i.repository.ListPaymentMethods(ctx)
}
