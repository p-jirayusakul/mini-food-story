package exceptions

import (
	"fmt"
	"strconv"
)

type Code string

const (
	CodeDomain     Code = "10000"
	CodeBusiness   Code = "11000"
	CodeSystem     Code = "12000"
	CodeRepository Code = "13000"

	CodeNotFound            Code = "14000"
	CodeProductNotFound     Code = "14001"
	CodeTableNotFound       Code = "14002"
	CodeTableStatusNotFound Code = "14003"
	CodeOrderNotFound       Code = "14004"
	CodeOrderItemNotFound   Code = "14005"
	CodeOrderStatusNotFound Code = "14006"
	CodeSessionFound        Code = "14007"

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
	ID      int64
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

func ErrorIDNotFound(code Code, id int64) *AppError {
	return &AppError{
		Code:    code,
		ID:      id,
		Message: notFoundMapping(id, code),
	}
}

func ErrorSessionNotFound() *AppError {
	return &AppError{
		Code:    CodeSessionFound,
		Message: "session not found",
	}
}

func ErrorDataNotFound() *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: "data not found",
	}
}

func Errorf(code Code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func notFoundMapping(id int64, code Code) string {
	u64, _ := strconv.ParseUint(string(code), 10, 16)
	codeInt := uint16(u64)
	if codeInt >= 14000 && codeInt < 14999 {
		var title string
		switch code {
		case CodeProductNotFound:
			title = fmt.Sprintf("product id '%d' not found", id)
			if id == 0 {
				title = "product not found"
			}
		case CodeTableNotFound:
			title = fmt.Sprintf("table id '%d' not found", id)
			if id == 0 {
				title = "table not found"
			}
		case CodeTableStatusNotFound:
			title = fmt.Sprintf("table status id '%d' not found", id)
			if id == 0 {
				title = "table status not found"
			}
		case CodeOrderNotFound:
			title = fmt.Sprintf("order id '%d' not found", id)
			if id == 0 {
				title = "order not found"
			}
		case CodeOrderItemNotFound:
			title = fmt.Sprintf("order item id '%d' not found", id)
			if id == 0 {
				title = "order item not found"
			}
		case CodeOrderStatusNotFound:
			title = fmt.Sprintf("order status id '%d' not found", id)
			if id == 0 {
				title = "order status not found"
			}
		case CodeSessionFound:
			title = "session not found"
		default:
			title = "data not found"
		}
		return title
	}
	return ""
}
