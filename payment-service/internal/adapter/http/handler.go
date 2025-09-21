package http

import (
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// CreatePaymentTransaction godoc
// @Summary Create payment transaction
// @Description Create a new payment transaction for an order
// @Tags Payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payment body Payment true "Payment transaction details"
// @Success 201 {object} middleware.SuccessResponse{data=createResponse}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router / [post]
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

// CallbackPaymentTransaction godoc
// @Summary Handle payment transaction callback
// @Description Process callback for payment transaction
// @Tags Payment
// @Accept json
// @Produce json
// @Param callback body CallbackPayment true "Payment callback details"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /callback [post]
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

	return middleware.ResponseOK(c, "payment callback processed successfully", nil)
}

// ListPaymentMethods godoc
// @Summary List payment methods
// @Description Get list of available payment methods
// @Tags Payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} middleware.SuccessResponse{data=[]domain.PaymentMethod}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /methods [get]
func (s *Handler) ListPaymentMethods(c *fiber.Ctx) error {
	result, customError := s.useCase.ListPaymentMethods(c.Context())
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get list payment methods success", result)
}

// GetPaymentLastStatusCodeByTransaction godoc
// @Summary List payment methods
// @Description Get list of available payment methods
// @Tags Payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} middleware.SuccessResponse{data=[]domain.PaymentMethod}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /methods [get]
func (s *Handler) GetPaymentLastStatusCodeByTransaction(c *fiber.Ctx) error {

	transactionID := c.Params("transactionID")
	result, customError := s.useCase.GetPaymentLastStatusCodeByTransaction(c.Context(), transactionID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	response := LastStatusCodeResponse{
		Code: result,
	}

	return middleware.ResponseOK(c, "get payment last status", response)
}
