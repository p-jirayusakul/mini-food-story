package http

import (
	"food-story/order/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strconv"
)

func (s *Handler) CreateOrder(c *fiber.Ctx) error {

	sessionIDData := c.Get("X-Session-Id")
	if sessionIDData == "" {
		return middleware.ResponseError(fiber.StatusBadRequest, "session id is required")
	}

	sessionID, err := uuid.Parse(sessionIDData)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.CreateOrder(c.Context(), sessionID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseCreated(c, "create order success", createResponse{
		ID: strconv.FormatInt(result, 10),
	})
}

func (s *Handler) GetOrderByID(c *fiber.Ctx) error {
	id, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetOrderByID(c.Context(), id)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order success", result)
}

func (s *Handler) UpdateOrderStatus(c *fiber.Ctx) error {

	sessionIDData := c.Get("X-Session-Id")
	if sessionIDData == "" {
		return middleware.ResponseError(fiber.StatusBadRequest, "session id is required")
	}

	sessionID, err := uuid.Parse(sessionIDData)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	id, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	body := new(Status)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	customError := s.useCase.UpdateOrderStatus(c.Context(), sessionID, domain.OrderStatus{
		ID:         id,
		StatusCode: body.StatusCode,
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "update order status success", nil)
}

func (s *Handler) CreateOrderItems(c *fiber.Ctx) error {
	sessionIDData := c.Get("X-Session-Id")
	if sessionIDData == "" {
		return middleware.ResponseError(fiber.StatusBadRequest, "session id is required")
	}

	sessionID, err := uuid.Parse(sessionIDData)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	id, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	body := new(OrderItems)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	var items []domain.OrderItems
	for _, item := range body.Items {
		productID, err := utils.StrToInt64(item.ProductID)
		if err != nil {
			return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
		}
		items = append(items, domain.OrderItems{
			OrderID:   id,
			ProductID: productID,
			Quantity:  item.Quantity,
			Note:      item.Note,
		})
	}
	customError := s.useCase.CreateOrderItems(c.Context(), sessionID, items)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseCreated(c, "create order item success", nil)
}

func (s *Handler) GetOrderItems(c *fiber.Ctx) error {
	id, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetOrderItems(c.Context(), id)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order items success", result)
}

func (s *Handler) GetOrderItemsByID(c *fiber.Ctx) error {
	orderID, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	orderItemsID, err := utils.StrToInt64(c.Params("orderItemsID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetOderItemsByID(c.Context(), orderID, orderItemsID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order item success", result)
}

func (s *Handler) UpdateOrderItemsStatus(c *fiber.Ctx) error {
	sessionIDData := c.Get("X-Session-Id")
	if sessionIDData == "" {
		return middleware.ResponseError(fiber.StatusBadRequest, "session id is required")
	}

	sessionID, err := uuid.Parse(sessionIDData)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	id, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	orderItemsID, err := utils.StrToInt64(c.Params("orderItemsID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	body := new(Status)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	customError := s.useCase.UpdateOrderItemsStatus(c.Context(), sessionID, domain.OrderItemsStatus{
		OrderID:    id,
		ID:         orderItemsID,
		StatusCode: body.StatusCode,
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "update order item status success", nil)
}
