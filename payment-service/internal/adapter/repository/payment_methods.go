package repository

import (
	"context"
	"fmt"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
)

func (i *PaymentRepoImplement) ListPaymentMethods(ctx context.Context) (result []*domain.PaymentMethod, customError *exceptions.CustomError) {
	data, err := i.repository.ListPaymentMethods(ctx)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch table status: %w", err),
		}
	}

	if data == nil {
		return nil, nil
	}

	result = make([]*domain.PaymentMethod, len(data))
	for index, v := range data {
		result[index] = &domain.PaymentMethod{
			ID:     v.ID,
			Code:   v.Code,
			Name:   v.Name,
			NameEN: v.NameEN,
		}
	}

	return result, nil
}
