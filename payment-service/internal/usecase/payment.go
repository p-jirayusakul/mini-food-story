package usecase

import (
	"context"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
)

func (i *PaymentImplement) CreatePaymentTransaction(ctx context.Context, payload domain.Payment) (transactionID string, customError *exceptions.CustomError) {

	customError = i.repository.IsOrderExist(ctx, payload.OrderID)
	if customError != nil {
		return "", customError
	}

	customError = i.repository.IsOrderItemsNotFinal(ctx, payload.OrderID)
	if customError != nil {
		return "", customError
	}

	transactionID, customError = i.repository.CreatePaymentTransaction(ctx, payload)
	if customError != nil {
		return "", customError
	}

	return transactionID, nil
}

func (i *PaymentImplement) GetPaymentLastStatusCodeByTransaction(ctx context.Context, transactionID string) (result string, customError *exceptions.CustomError) {
	return i.repository.GetPaymentLastStatusCodeByTransaction(ctx, transactionID)
}

func (i *PaymentImplement) CallbackPaymentTransaction(ctx context.Context, transactionID string) (customError *exceptions.CustomError) {

	customError = i.repository.CallbackPaymentTransaction(ctx, transactionID)
	if customError != nil {
		return customError
	}

	return nil
}
