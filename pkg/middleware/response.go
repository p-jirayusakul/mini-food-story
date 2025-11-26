package middleware

import (
	"errors"
	"food-story/pkg/exceptions"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SuccessResponse struct {
	Meta successDetail `json:"meta"`
	Data interface{}   `json:"data" swaggertype:"object"`
}

type SuccessWithPaginationResponse struct {
	Meta successWithPaginationDetail `json:"meta"`
	Data interface{}                 `json:"data" swaggertype:"object"`
}

type successDetail struct {
	RequestId string `json:"requestId"`
	Timestamp string `json:"timestamp"`
}

type successWithPaginationDetail struct {
	PageSize   int64  `json:"pageSize"`
	PageNumber int64  `json:"pageNumber"`
	TotalItems int64  `json:"totalItems"`
	TotalPages int64  `json:"totalPages"`
	RequestId  string `json:"requestId"`
	Timestamp  string `json:"timestamp"`
}

type ErrorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceId string `json:"traceId"`
}

type ResponseWithPaginationPayload struct {
	PageSize   int64
	PageNumber int64
	TotalItems int64
	TotalPages int64
	Data       interface{}
}

func ResponseOK(c *fiber.Ctx, payload interface{}) error {

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{Meta: successDetail{RequestId: getRequestID(c), Timestamp: getTimestamp()}, Data: payload})
}

func ResponseOKWithPagination(c *fiber.Ctx, payload ResponseWithPaginationPayload) error {
	return c.Status(fiber.StatusOK).JSON(SuccessWithPaginationResponse{Meta: successWithPaginationDetail{PageSize: payload.PageSize, PageNumber: payload.PageNumber, TotalItems: payload.TotalItems, TotalPages: payload.TotalPages, RequestId: getRequestID(c), Timestamp: getTimestamp()}, Data: payload.Data})
}

func ResponseCreated(c *fiber.Ctx, payload interface{}) error {

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{Meta: successDetail{RequestId: getRequestID(c), Timestamp: getTimestamp()}, Data: payload})
}

func ResponseError(c *fiber.Ctx, err error) error {
	httpCode, response := mapErrorToHTTP(err)
	if httpCode >= 500 {
		slog.Error(err.Error())
	}
	response.Error.TraceId = getRequestID(c)
	return c.Status(httpCode).JSON(response)
}

func mapErrorToHTTP(err error) (int, ErrorResponse) {
	var appErr *exceptions.AppError

	if errors.As(err, &appErr) {

		switch appErr.Code {

		case exceptions.CodeDomain:
			return 400, getErrorResponse(string(appErr.Code), appErr.Message)

		case exceptions.CodeBusiness:
			return 400, getErrorResponse(string(appErr.Code), appErr.Message)

		case exceptions.CodeUnauthorized:
			return 401, getErrorResponse(string(appErr.Code), appErr.Message)

		case exceptions.CodeForbidden:
			return 403, getErrorResponse(string(appErr.Code), appErr.Message)

		case exceptions.CodeNotFound, exceptions.CodeOrderNotFound, exceptions.CodeTableNotFound, exceptions.CodeOrderItemNotFound, exceptions.CodeProductNotFound, exceptions.CodeTableStatusNotFound, exceptions.CodeSessionFound:
			return 404, getErrorResponse(string(appErr.Code), appErr.Message)

		case exceptions.CodeConflict:
			return 409, getErrorResponse(string(appErr.Code), appErr.Message)

		case exceptions.CodeRepository,
			exceptions.CodeRedis,
			exceptions.CodeSystem:
			return 500, getErrorResponse(string(appErr.Code), appErr.Message)
		}
	}

	return 500, getErrorResponse(string(exceptions.CodeUnknown), "unknown error")
}

func getErrorResponse(code string, message string) ErrorResponse {
	return ErrorResponse{Error: errorDetail{Code: code, Message: message}}
}

func getTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func getRequestID(c *fiber.Ctx) string {
	reqID, ok := c.Locals(CtxRequestIDKey).(string)
	if !ok {
		reqID = ""
	}
	return reqID
}
