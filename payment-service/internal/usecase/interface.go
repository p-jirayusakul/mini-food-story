package usecase

import (
	"context"
	"food-story/payment-service/internal/adapter/repository"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/shared/config"
)

type PaymentUsecase interface {
	ListPaymentMethods(ctx context.Context) (result []*domain.PaymentMethod, customError *exceptions.CustomError)
	CreatePaymentTransaction(ctx context.Context, payload domain.Payment) (transactionID string, customError *exceptions.CustomError)
	CallbackPaymentTransaction(ctx context.Context, transactionID string) (customError *exceptions.CustomError)
}

type PaymentImplement struct {
	config     config.Config
	repository repository.Implement
}

func NewUsecase(config config.Config, repository repository.Implement) *PaymentImplement {
	return &PaymentImplement{
		config,
		repository,
	}
}

var _ PaymentUsecase = (*PaymentImplement)(nil)
