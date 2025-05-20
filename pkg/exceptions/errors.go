package exceptions

import (
	"errors"
	"github.com/jackc/pgx/v5"
)

var (
	ErrIDInvalidFormat      = errors.New("invalid id format")
	ErrValueIsEmpty         = errors.New("value is empty")
	ErrInternalServerError  = errors.New("something went wrong please try again")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrRowDatabaseNotFound  = pgx.ErrNoRows
	ErrRedisKeyNotFound     = errors.New("key not found")
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionExpired       = errors.New("session expired")
	ErrFailedToReadSession  = errors.New("failed to read session")
	ErrCtxCanceledOrTimeout = errors.New("request cancelled or timeout")
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderItemsNotFound   = errors.New("order items not found")
)

const (
	ErrMsgDataNotFound = "data not found"
)

const (
	SqlstateUniqueViolation = "23505"
)
