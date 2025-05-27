package http

import (
	"food-story/kitchen-service/internal/domain"
	"food-story/pkg/common"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"
	"github.com/gofiber/fiber/v2"
)

const ResGetOrderItemsMsg = "get order items success"

// SearchOrderItems godoc
// @Summary Search order items
// @Description Search order items with filters
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param pageNumber query int false "Page number" minimum(1)
// @Param pageSize query int false "Page size" minimum(1)
// @Param search query string false "Search by name" maxLength(255)
// @Param statusCode query []string false "Filter by status codes" Enums(PENDING, PROCESSING, SERVED, CANCELLED)
// @Param tableNumber query []string false "Filter by table numbers"
// @Param orderBy query string false "Order by field" Enums(id,tableNumber,statusCode,productName,quantity)
// @Param orderType query string false "Order direction" Enums(asc,desc)
// @Success 200 {object} middleware.SuccessResponse{data=domain.SearchOrderItemsResult}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders/search/items [get]
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

	return middleware.ResponseOK(c, ResGetOrderItemsMsg, result)
}

// GetOrderItems godoc
// @Summary Get order items for specific order
// @Description Get order items by order ID with pagination
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param pageNumber query int false "Page number" minimum(1)
// @Success 200 {object} middleware.SuccessResponse{data=domain.SearchOrderItemsResult}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders/{id}/items [get]
func (s *Handler) GetOrderItems(c *fiber.Ctx) error {
	orderID, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	body := new(SearchOrderItemsByOrderID)
	if err := c.QueryParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	payload := domain.SearchOrderItems{
		PageSize:   common.DefaultPageSize,
		PageNumber: body.PageNumber,
	}

	result, customError := s.useCase.GetOrderItems(c.Context(), orderID, payload)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, ResGetOrderItemsMsg, result)
}

// GetOrderItemsByID godoc
// @Summary Get specific order item
// @Description Get order item by order ID and order item ID
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param orderItemsID path int true "Order Item ID"
// @Success 200 {object} middleware.SuccessResponse{data=model.OrderItems}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders/{id}/items/{orderItemsID} [get]
func (s *Handler) GetOrderItemsByID(c *fiber.Ctx) error {
	orderItemsID, orderID, err := handleParams(c)
	if err != nil {
		return err
	}

	result, customError := s.useCase.GetOrderItemsByID(c.Context(), orderID, orderItemsID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, ResGetOrderItemsMsg, result)
}

// UpdateOrderItemsStatusServe godoc
// @Summary Update order item status to serv
// @Description Update status of specific order item to serv
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param orderItemsID path int true "Order Item ID"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders/{id}/items/{orderItemsID}/status/serve [patch]
func (s *Handler) UpdateOrderItemsStatusServe(c *fiber.Ctx) error {
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

// UpdateOrderItemsStatusCancel godoc
// @Summary Update order item status to cancel
// @Description Update status of specific order item to cancel
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param orderItemsID path int true "Order Item ID"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /orders/{id}/items/{orderItemsID}/status/cancel [patch]
func (s *Handler) UpdateOrderItemsStatusCancel(c *fiber.Ctx) error {
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
