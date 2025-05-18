package http

import (
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func (s *Handler) CreatePaymentTransaction(c *fiber.Ctx) error {
	body := new(Payment)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	orderID, err := utils.StrToInt64(body.OrderID)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	method, err := utils.StrToInt64(body.Method)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.CreatePaymentTransaction(c.Context(), domain.Payment{
		OrderID: orderID,
		Method:  method,
		Note:    body.Note,
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseCreated(c, "create payment transaction success", createResponse{
		TransactionID: result,
	})
}

func (s *Handler) CallbackPaymentTransaction(c *fiber.Ctx) error {
	body := new(CallbackPayment)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	customError := s.useCase.CallbackPaymentTransaction(c.Context(), body.TransactionID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get callback success", nil)
}

func (s *Handler) ListPaymentMethods(c *fiber.Ctx) error {
	result, customError := s.useCase.ListPaymentMethods(c.Context())
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get list payment methods success", result)
}
