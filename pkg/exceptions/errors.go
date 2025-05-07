package exceptions

import (
	"errors"
	"github.com/jackc/pgx/v5"
)

var (
	ErrIDInvalidFormat     = errors.New("invalid id format")
	ErrValueIsEmpty        = errors.New("value is empty")
	ErrInternalServerError = errors.New("something went wrong please try again")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrRowDatabaseNotFound = pgx.ErrNoRows
)

const (
	EXCOK                  = "OK"
	EXCInternalServerError = "Internal Server Error"
	EXCDataNotFound        = "Data not found"
	EXCBadRequest          = "Bad Request"
	EXCConvertValue        = "error converting value"
	EXCOneOf               = "failed on the 'oneof' tag"
	EXCMin                 = "failed on the 'min' tag"
	EXCMax                 = "failed on the 'max' tag"
	EXCRequired            = "failed on the 'required' tag"
	EXCNoSpecialChar       = "failed on the 'no_special_char'"
	EXCNoSpecialCharSlash  = "failed on the 'no_special_char_slash'"
)

const (
	SqlstateUniqueViolation = "23505"
)
