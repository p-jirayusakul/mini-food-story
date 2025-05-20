package http

import (
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"
	"github.com/gofiber/fiber/v2"
)

func (s *Handler) SearchOrderItems(c *fiber.Ctx) error {
	body := new(SearchOrderItems)
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
		TableNumber: utils.FilterOutZeroInt(body.TableNumber),
		StatusCode:  utils.FilterOutEmptyStr(body.StatusCode),
		OrderByType: orderByType,
		OrderBy:     body.OrderBy,
		PageSize:    body.PageSize,
		PageNumber:  body.PageNumber,
	}

	result, customError := s.useCase.SearchOrderItems(c.Context(), payload)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order items success", result)
}

func (s *Handler) GetOrderItems(c *fiber.Ctx) error {
	orderID, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetOrderItems(c.Context(), orderID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order items success", result)
}

func (s *Handler) GetOrderItemsByID(c *fiber.Ctx) error {
	orderItemsID, orderID, err := handleParams(c)
	if err != nil {
		return err
	}

	result, customError := s.useCase.GetOrderItemsByID(c.Context(), orderID, orderItemsID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order items success", result)
}

func (s *Handler) UpdateOrderItemsStatusServed(c *fiber.Ctx) error {
	orderItemsID, orderID, err := handleParams(c)
	if err != nil {
		return err
	}

	customError := s.useCase.UpdateOrderItemsStatusServed(c.Context(), shareModel.OrderItemsStatus{
		ID:         orderItemsID,
		OrderID:    orderID,
		StatusCode: "SERVED",
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "update order item status served success", nil)
}

func (s *Handler) UpdateOrderItemsStatusCancelled(c *fiber.Ctx) error {
	orderItemsID, orderID, err := handleParams(c)
	if err != nil {
		return err
	}

	customError := s.useCase.UpdateOrderItemsStatus(c.Context(), shareModel.OrderItemsStatus{
		ID:         orderItemsID,
		OrderID:    orderID,
		StatusCode: "CANCELLED",
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "update order item status cancelled success", nil)
}

func handleParams(c *fiber.Ctx) (orderItemsID, orderID int64, err error) {
	orderItemsID, err = utils.StrToInt64(c.Params("orderItemsID"))
	if err != nil {
		return 0, 0, middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	orderID, err = utils.StrToInt64(c.Params("id"))
	if err != nil {
		return 0, 0, middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	return orderItemsID, orderID, nil
}
