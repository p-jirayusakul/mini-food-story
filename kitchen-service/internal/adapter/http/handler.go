package http

import (
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func (s *Handler) UpdateOrderItemsStatus(c *fiber.Ctx) error {
	orderItemsID, err := utils.StrToInt64(c.Params("orderItemsID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	orderID, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	body := new(StatusOrderItems)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	customError := s.useCase.UpdateOrderItemsStatus(c.Context(), domain.OrderItemsStatus{
		ID:         orderItemsID,
		OrderID:    orderID,
		StatusCode: body.StatusCode,
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "update order item status success", nil)
}
