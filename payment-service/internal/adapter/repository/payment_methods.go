package repository

import (
	"context"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
)

func (i *Implement) ListPaymentMethods(ctx context.Context) (result []*domain.PaymentMethod, err error) {
	data, err := i.repository.ListPaymentMethods(ctx)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch table status", err)
	}

	if data == nil {
		return nil, exceptions.ErrorDataNotFound()
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
