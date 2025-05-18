package usecase

import (
	"context"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
)

func (i *PaymentImplement) ListPaymentMethods(ctx context.Context) (result []*domain.PaymentMethod, customError *exceptions.CustomError) {
	return i.repository.ListPaymentMethods(ctx)
}
