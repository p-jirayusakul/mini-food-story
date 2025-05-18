package usecase

import (
	"context"
	"fmt"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
	"time"
)

func (i *PaymentImplement) CreatePaymentTransaction(ctx context.Context, payload domain.Payment) (transactionID string, customError *exceptions.CustomError) {

	customError = i.repository.IsOrderItemsNotFinal(ctx, payload.OrderID)
	if customError != nil {
		return "", customError
	}

	transactionID, customError = i.repository.CreatePaymentTransaction(ctx, payload)
	if customError != nil {
		return "", customError
	}

	customError = i.CallbackPaymentTransaction(ctx, transactionID)
	if customError != nil {
		return "", customError
	}

	return transactionID, nil
}

func (i *PaymentImplement) CallbackPaymentTransaction(ctx context.Context, transactionID string) (customError *exceptions.CustomError) {

	fmt.Println("Start call back payment transaction")
	time.Sleep(2 * time.Second)
	customError = i.repository.CallbackPaymentTransaction(ctx, transactionID)
	fmt.Println("End call back payment transaction")
	if customError != nil {
		return customError
	}

	return nil
}
