package http

import (
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func (s *Handler) CreateOrder(c *fiber.Ctx) error {

	sessionIDData := c.Get("X-Session-Id")
	if sessionIDData == "" {
		return middleware.ResponseError(fiber.StatusBadRequest, "session id is required")
	}

	sessionID, err := utils.DecryptSessionToUUID(sessionIDData, []byte(s.config.SecretKey))
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
			ProductID: productID,
			Quantity:  item.Quantity,
			Note:      item.Note,
		})
	}

	_, customError := s.useCase.CreateOrder(c.Context(), sessionID, items)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseCreated(c, "create order success", nil)
}

func (s *Handler) GetOrderByID(c *fiber.Ctx) error {
	sessionIDData := c.Get("X-Session-Id")
	if sessionIDData == "" {
		return middleware.ResponseError(fiber.StatusBadRequest, "session id is required")
	}
	sessionID, err := utils.DecryptSessionToUUID(sessionIDData, []byte(s.config.SecretKey))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetOrderByID(c.Context(), sessionID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order success", result)
}

func (s *Handler) CreateOrderItems(c *fiber.Ctx) error {
	sessionIDData := c.Get("X-Session-Id")
	if sessionIDData == "" {
		return middleware.ResponseError(fiber.StatusBadRequest, "session id is required")
	}

	sessionID, err := utils.DecryptSessionToUUID(sessionIDData, []byte(s.config.SecretKey))
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
	sessionIDData := c.Get("X-Session-Id")
	if sessionIDData == "" {
		return middleware.ResponseError(fiber.StatusBadRequest, "session id is required")
	}
	sessionID, err := utils.DecryptSessionToUUID(sessionIDData, []byte(s.config.SecretKey))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetOrderItems(c.Context(), sessionID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order items success", result)
}

func (s *Handler) GetOrderItemsByID(c *fiber.Ctx) error {
	sessionIDData := c.Get("X-Session-Id")
	if sessionIDData == "" {
		return middleware.ResponseError(fiber.StatusBadRequest, "session id is required")
	}
	sessionID, err := utils.DecryptSessionToUUID(sessionIDData, []byte(s.config.SecretKey))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	orderItemsID, err := utils.StrToInt64(c.Params("orderItemsID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetOderItemsByID(c.Context(), sessionID, orderItemsID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order item success", result)
}

func (s *Handler) UpdateOrderItemsStatusCancelled(c *fiber.Ctx) error {
	sessionIDData := c.Get("X-Session-Id")
	if sessionIDData == "" {
		return middleware.ResponseError(fiber.StatusBadRequest, "session id is required")
	}

	sessionID, err := utils.DecryptSessionToUUID(sessionIDData, []byte(s.config.SecretKey))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	orderItemsID, err := utils.StrToInt64(c.Params("orderItemsID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	customError := s.useCase.UpdateOrderItemsStatus(c.Context(), sessionID, domain.OrderItemsStatus{
		ID:         orderItemsID,
		StatusCode: "CANCELLED",
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "update order item status success", nil)
}

func (s *Handler) SearchOrderItemsInComplete(c *fiber.Ctx) error {
	orderID, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	body := new(SearchOrderItemsIncomplete)
	if err := c.QueryParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	orderByType := "desc"
	if body.OrderByType != "" {
		orderByType = body.OrderByType
	}

	payload := domain.SearchOrderItems{
		Name:        body.Search,
		StatusCode:  utils.FilterOutEmptyStr(body.StatusCode),
		OrderByType: orderByType,
		OrderBy:     body.OrderBy,
		PageSize:    body.PageSize,
		PageNumber:  body.PageNumber,
	}

	result, customError := s.useCase.SearchOrderItemsIncomplete(c.Context(), orderID, payload)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order items success", result)
}
