package http

import (
	"context"
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func (s *Handler) SearchOrderItems(c *fiber.Ctx) error {
	result, customError := s.useCase.SearchOrderItems(c.Context())
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

	err = s.HandleStatusOrderItems(c.Context(), domain.OrderItemsStatus{
		ID:         orderItemsID,
		OrderID:    orderID,
		StatusCode: "SERVED",
	})
	if err != nil {
		return err
	}

	return middleware.ResponseOK(c, "update order item status served success", nil)
}

func (s *Handler) UpdateOrderItemsStatusCancelled(c *fiber.Ctx) error {
	orderItemsID, orderID, err := handleParams(c)
	if err != nil {
		return err
	}

	err = s.HandleStatusOrderItems(c.Context(), domain.OrderItemsStatus{
		ID:         orderItemsID,
		OrderID:    orderID,
		StatusCode: "CANCELLED",
	})
	if err != nil {
		return err
	}

	return middleware.ResponseOK(c, "update order item status cancelled success", nil)
}

func (s *Handler) HandleStatusOrderItems(ctx context.Context, payload domain.OrderItemsStatus) error {
	customError := s.useCase.UpdateOrderItemsStatus(ctx, payload)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return nil
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
