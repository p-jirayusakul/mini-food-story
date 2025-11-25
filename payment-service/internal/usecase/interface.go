package usecase

import (
	"context"
	"food-story/payment-service/internal/adapter/cache"
	"food-story/payment-service/internal/adapter/repository"
	"food-story/payment-service/internal/domain"
	"food-story/shared/config"
)

type PaymentUsecase interface {
	ListPaymentMethods(ctx context.Context) (result []*domain.PaymentMethod, err error)
	CreatePaymentTransaction(ctx context.Context, payload domain.Payment) (transactionID string, err error)
	CallbackPaymentTransaction(ctx context.Context, transactionID string, statusCode string) (err error)
	GetPaymentLastStatusCodeByTransaction(ctx context.Context, transactionID string) (result string, err error)
	PaymentTransactionQR(ctx context.Context, transactionID string) (result domain.TransactionQR, err error)
}

type PaymentImplement struct {
	config     config.Config
	repository repository.Implement
	cache      cache.RedisTableCacheInterface
}

func NewUsecase(config config.Config, repository repository.Implement, cache cache.RedisTableCacheInterface) *PaymentImplement {
	return &PaymentImplement{
		config,
		repository,
		cache,
	}
}

var _ PaymentUsecase = (*PaymentImplement)(nil)
