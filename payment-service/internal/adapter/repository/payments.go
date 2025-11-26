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
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const FailToGetTotalAmount = "failed to fetch amount order: %w"

func (i *Implement) CreatePaymentTransaction(ctx context.Context, payload domain.Payment) (transactionID string, err error) {

	var note pgtype.Text
	if payload.Note != nil {
		note.Valid = true
		note.String = *payload.Note
	}

	amount, err := i.repository.GetTotalAmountToPayForServedItems(ctx, payload.OrderID)
	if err != nil {
		return "", exceptions.Errorf(exceptions.CodeRepository, FailToGetTotalAmount, err)
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
		return "", exceptions.Errorf(exceptions.CodeRepository, "failed to create payment", err)
	}

	err = i.repository.UpdateOrderStatusWaitForPayment(ctx, payload.OrderID)
	if err != nil {
		return "", exceptions.Errorf(exceptions.CodeRepository, "failed to update order status", err)
	}

	tableID, customError := i.GetTableIDByOrderID(ctx, payload.OrderID)
	if customError != nil {
		return "", customError
	}

	customError = i.UpdateTablesStatusWaitingForPayment(ctx, tableID)
	if customError != nil {
		return "", customError
	}

	return transactionID, nil
}
func (i *Implement) GetPaymentLastStatusCodeByTransaction(ctx context.Context, transactionID string) (statusCode string, err error) {

	const errNotFound = "payment status not found"
	result, err := i.repository.GetPaymentLastStatusCodeByTransaction(ctx, transactionID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return "", exceptions.Error(exceptions.CodeNotFound, errNotFound)
		}
		return "", exceptions.Errorf(exceptions.CodeRepository, "failed to get payment status", err)
	}

	if result.String == "" {
		return "", exceptions.Error(exceptions.CodeRepository, errNotFound)
	}

	return result.String, nil
}

func (i *Implement) CallbackPaymentTransaction(ctx context.Context, transactionID, statusCode string) (sessionID uuid.UUID, err error) {

	if statusCode == "" || statusCode == "PENDING" {
		return uuid.Nil, nil
	}

	if strings.ToUpper(statusCode) != "SUCCESS" {
		switch statusCode {
		case "CANCELLED":
			err := i.repository.UpdateStatusPaymentCancelledByTransactionID(ctx, transactionID)
			if err != nil {
				return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, "failed to update payment status cancelled", err)
			}
		case "FAILED":
			err := i.repository.UpdateStatusPaymentFailedByTransactionID(ctx, transactionID)
			if err != nil {
				return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, "failed to update payment status failed", err)
			}
		case "TIMEOUT":
			err := i.repository.UpdateStatusPaymentTimeOutByTransactionID(ctx, transactionID)
			if err != nil {
				return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, "failed to update payment status timeout", err)
			}
		default:
			return uuid.Nil, exceptions.Error(exceptions.CodeRepository, "status code not supported")
		}
	} else {
		err := i.repository.UpdateStatusPaymentSuccessByTransactionID(ctx, transactionID)
		if err != nil {
			return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, "failed to update payment status success", err)
		}
	}

	orderID, err := i.repository.GetPaymentOrderIDByTransaction(ctx, transactionID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return uuid.Nil, exceptions.ErrorIDNotFound(exceptions.CodeOrderNotFound, 0)
		}

		return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, "failed to get order id", err)
	}

	if strings.ToUpper(statusCode) != "SUCCESS" {
		err = i.repository.UpdateOrderStatusWaitForPayment(ctx, orderID)
		if err != nil {
			return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, "failed to update order status", err)
		}
		return uuid.Nil, nil
	}

	amount, err := i.repository.GetTotalAmountToPayForServedItems(ctx, orderID)
	if err != nil {
		return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, FailToGetTotalAmount, err)
	}

	err = i.repository.UpdateOrderStatusCompletedAndAmount(ctx, database.UpdateOrderStatusCompletedAndAmountParams{
		ID:     orderID,
		Amount: amount,
	})
	if err != nil {
		return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, "failed to update order status and amount", err)
	}

	tableID, customError := i.GetTableIDByOrderID(ctx, orderID)
	if customError != nil {
		return uuid.Nil, customError
	}

	customError = i.UpdateTablesStatusCleaning(ctx, tableID)
	if customError != nil {
		return uuid.Nil, customError
	}

	sessionID, customError = i.GetSessionIDByOrderID(ctx, orderID)
	if customError != nil {
		return uuid.Nil, customError
	}

	customError = i.UpdateStatusCloseTableSession(ctx, sessionID)
	if customError != nil {
		return uuid.Nil, customError
	}

	return sessionID, nil
}

func (i *Implement) GetPaymentAmountByTransaction(ctx context.Context, transactionID string) (result float64, err error) {
	amount, err := i.repository.GetPaymentAmountByTransaction(ctx, transactionID)
	if err != nil {
		if errors.Is(err, exceptions.ErrRowDatabaseNotFound) {
			return 0, exceptions.Errorf(exceptions.CodeNotFound, FailToGetTotalAmount, err)
		}
		return 0, exceptions.Errorf(exceptions.CodeRepository, "failed to get payment amount", err)
	}
	amountFloat, err := amount.Float64Value()
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeSystem, "failed to convert amount to float64", err)
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
