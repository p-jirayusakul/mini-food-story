package repository

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
	"log/slog"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const FailToGetTotalAmount = "failed to fetch amount order: %w"

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
			Errors: fmt.Errorf(FailToGetTotalAmount, err),
		}
	}

	transactionID = uuid.New().String()
	pendingStatus, customError := i.GetPaymentStatusPending(ctx)
	if customError != nil {
		return "", customError
	}
	arg := database.CreatePaymentParams{
		ID:            i.snowflakeID.Generate(),
		OrderID:       payload.OrderID,
		Method:        payload.Method,
		Status:        pendingStatus,
		Note:          note,
		Amount:        amount,
		TransactionID: transactionID,
		RefCode:       GenerateRefCode(),
	}

	_, err = i.repository.CreatePayment(ctx, arg)
	if err != nil {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to create payment: %w", err),
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

func (i *Implement) GetPaymentLastStatusCodeByTransaction(ctx context.Context, transactionID string) (statusCode string, customError *exceptions.CustomError) {

	const errNotFound = "payment status not found"
	result, err := i.repository.GetPaymentLastStatusCodeByTransaction(ctx, transactionID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return "", &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: fmt.Errorf(errNotFound),
			}
		}
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}

	if result.String == "" {
		return "", &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf(errNotFound),
		}
	}

	return result.String, nil
}

func (i *Implement) CallbackPaymentTransaction(ctx context.Context, transactionID string) (customError *exceptions.CustomError) {
	err := i.repository.UpdateStatusPaymentSuccessByTransactionID(ctx, transactionID)
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
			Errors: fmt.Errorf(FailToGetTotalAmount, err),
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

	tableID, customError := i.GetTableIDByOrderID(ctx, orderID)
	if customError != nil {
		return customError
	}

	customError = i.UpdateTablesStatusCleaning(ctx, tableID)
	if customError != nil {
		return customError
	}

	return nil
}

func (i *Implement) GetPaymentAmountByTransaction(ctx context.Context, transactionID string) (result float64, customError *exceptions.CustomError) {
	amount, err := i.repository.GetPaymentAmountByTransaction(ctx, transactionID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, &exceptions.CustomError{
				Status: exceptions.ERRNOTFOUND,
				Errors: fmt.Errorf(FailToGetTotalAmount, err),
			}
		}
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: err,
		}
	}
	amountFloat, err := amount.Float64Value()
	if err != nil {
		return 0, &exceptions.CustomError{
			Status: exceptions.ERRSYSTEM,
			Errors: err,
		}
	}

	return amountFloat.Float64, nil
}

func GenerateRefCode() string {
	now := time.Now()
	datePart := now.Format("20060102")

	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		slog.Error("failed to generate secure random number: ", "err", err)
		return ""
	}
	randomPart := n.Int64() + 100000

	return fmt.Sprintf("PAY-%s-%d", datePart, randomPart)
}
