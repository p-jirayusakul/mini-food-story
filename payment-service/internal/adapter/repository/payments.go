package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"math/rand"
	"time"
)

func (i *Implement) CreatePaymentTransaction(ctx context.Context, payload domain.Payment) (transactionID string, customError *exceptions.CustomError) {

	var note pgtype.Text
	if payload.Note != nil {
		note.Valid = true
		note.String = *payload.Note
	}

	amount, err := i.repository.GetTotalAmountToPayForServedItems(ctx, payload.OrderID)
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch table status: %w", err),
		}
	}

	transactionID = uuid.New().String()
	arg := database.CreatePaymentParams{
		ID:            i.snowflakeID.Generate(),
		OrderID:       payload.OrderID,
		Method:        payload.Method,
		Note:          note,
		Amount:        amount,
		TransactionID: transactionID,
		RefCode:       GenerateRefCode(),
	}

	_, err = i.repository.CreatePayment(ctx, arg)
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch table status: %w", err),
		}
	}

	err = i.repository.UpdateOrderStatusWaitForPayment(ctx, payload.OrderID)
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
		}
	}

	return transactionID, nil
}

func (i *Implement) CallbackPaymentTransaction(ctx context.Context, transactionID string) (customError *exceptions.CustomError) {
	err := i.repository.UpdateStatusPaymentPaidByTransactionID(ctx, transactionID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	orderID, err := i.repository.GetPaymentOrderIDByTransaction(ctx, transactionID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: exceptions.ErrOrderNotFound,
			}
		}

		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	amount, err := i.repository.GetTotalAmountToPayForServedItems(ctx, orderID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch table status: %w", err),
		}
	}
	err = i.repository.UpdateOrderStatusCompletedAndAmount(ctx, database.UpdateOrderStatusCompletedAndAmountParams{
		ID:     orderID,
		Amount: amount,
	})
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	return nil
}

func GenerateRefCode() string {
	now := time.Now()
	datePart := now.Format("20060102")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomPart := r.Intn(900000) + 100000

	return fmt.Sprintf("PAY-%s-%d", datePart, randomPart)
}
