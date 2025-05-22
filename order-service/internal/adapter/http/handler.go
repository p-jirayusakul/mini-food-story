package http

import (
	"food-story/order-service/internal/domain"
	"food-story/pkg/common"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *Handler) CreateOrder(c *fiber.Ctx) error {

	sessionID, err := getSession(c)
	if err != nil {
		return err
	}

	body := new(OrderItems)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	var items []shareModel.OrderItems
	for _, item := range body.Items {
		productID, err := utils.StrToInt64(item.ProductID)
		if err != nil {
			return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
		}
		items = append(items, shareModel.OrderItems{
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
	sessionID, err := getSession(c)
	if err != nil {
		return err
	}

	result, customError := s.useCase.GetOrderByID(c.Context(), sessionID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order success", CurrentOrderResponse{
		TableNumber:  result.TableNumber,
		StatusName:   result.StatusName,
		StatusNameEN: result.StatusNameEN,
		StatusCode:   result.StatusCode,
	})
}

func (s *Handler) CreateOrderItems(c *fiber.Ctx) error {
	sessionID, err := getSession(c)
	if err != nil {
		return err
	}

	body := new(OrderItems)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	var items []shareModel.OrderItems
	for _, item := range body.Items {
		productID, err := utils.StrToInt64(item.ProductID)
		if err != nil {
			return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
		}
		items = append(items, shareModel.OrderItems{
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
	sessionID, err := getSession(c)
	if err != nil {
		return err
	}

	body := new(SearchOrderItems)
	if errValidate := c.QueryParser(body); errValidate != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, errValidate.Error())
	}

	if err = s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	payload := domain.SearchOrderItems{
		OrderBy:    body.OrderBy,
		PageSize:   int64(common.DefaultPageSize),
		PageNumber: body.PageNumber,
	}

	result, customError := s.useCase.GetCurrentOrderItems(c.Context(), sessionID, payload)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order items success", result)
}

func (s *Handler) GetOrderItemsByID(c *fiber.Ctx) error {
	sessionID, err := getSession(c)
	if err != nil {
		return err
	}

	orderItemsID, err := utils.StrToInt64(c.Params("orderItemsID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetCurrentOrderItemsByID(c.Context(), sessionID, orderItemsID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order item success", result)
}

func (s *Handler) UpdateOrderItemsStatusCancelled(c *fiber.Ctx) error {
	sessionID, err := getSession(c)
	if err != nil {
		return err
	}

	orderItemsID, err := utils.StrToInt64(c.Params("orderItemsID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	customError := s.useCase.UpdateOrderItemsStatus(c.Context(), sessionID, shareModel.OrderItemsStatus{
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

func getSession(c *fiber.Ctx) (uuid.UUID, error) {
	sessionIDAny, ok := c.Locals("sessionID").(string)
	if !ok {
		return uuid.UUID{}, middleware.ResponseError(fiber.StatusInternalServerError, exceptions.ErrFailedToReadSession.Error())
	}

	sessionID, err := utils.PareStringToUUID(sessionIDAny)
	if err != nil {
		return uuid.UUID{}, middleware.ResponseError(fiber.StatusInternalServerError, exceptions.ErrFailedToReadSession.Error())
	}

	return sessionID, nil
}
