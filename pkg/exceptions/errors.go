package exceptions

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrIDInvalidFormat      = errors.New("invalid id format")
	ErrValueIsEmpty         = errors.New("value is empty")
	ErrInternalServerError  = errors.New("something went wrong please try again")
	ErrRowDatabaseNotFound  = pgx.ErrNoRows
	ErrRedisKeyNotFound     = errors.New("key not found")
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionExpired       = errors.New("session expired")
	ErrFailedToReadSession  = errors.New("failed to read session")
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderRequired        = errors.New("order cannot be empty")
	ErrOrderItemsNotFound   = errors.New("order items not found")
	ErrOrderItemsRequired   = errors.New("order items cannot be empty")
	ErrProductNotFound      = errors.New("product not found")
	ErrTableNotFound        = errors.New("table not found")
	ErrTableSessionNotFound = errors.New("table session not found")

	ErrRedisKeyNotFoundException = Error(CodeNotFound, "key not found")
)
