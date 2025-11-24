package exceptions

import (
	"fmt"
)

type Code string

const (
	CodeDomain       Code = "10000"
	CodeBusiness     Code = "11000"
	CodeSystem       Code = "12000"
	CodeRepository   Code = "13000"
	CodeNotFound     Code = "14000"
	CodeUnauthorized Code = "15000"
	CodeForbidden    Code = "16000"
	CodeConflict     Code = "17000"
	CodeRedis        Code = "18000"
	CodeUnknown      Code = "90000"
)

type AppError struct {
	Code    Code
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s | root: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func Error(code Code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Errorf(code Code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
