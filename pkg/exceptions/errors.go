package exceptions

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrIDInvalidFormat     = errors.New("invalid id format")
	ErrValueIsEmpty        = errors.New("value is empty")
	ErrInternalServerError = errors.New("something went wrong please try again")
	ErrRowDatabaseNotFound = pgx.ErrNoRows
	ErrRedisKeyNotFound    = errors.New("key not found")
	ErrSessionExpired      = errors.New("session expired")
	ErrFailedToReadSession = errors.New("failed to read session")
	ErrOrderRequired       = errors.New("order cannot be empty")
	ErrOrderItemsRequired  = errors.New("order items cannot be empty")
	ErrForeignKeyViolation = errors.New("foreign key violation")

	ErrRedisKeyNotFoundException = Error(CodeNotFound, "key not found")
)
