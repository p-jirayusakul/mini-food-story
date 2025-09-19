package http

import (
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateOrder godoc
// @Summary Create new order
// @Description Create a new order with items for current table session
// @Tags Order
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Param order body OrderItems true "Order item details"
// @Success 201 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /current [post]
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

// GetOrderByID godoc
// @Summary Get order details by session ID
// @Description Get current order details for the given session ID
// @Tags Order
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Success 200 {object} CurrentOrderResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /current [get]
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

// CreateOrderItems godoc
// @Summary Add items to an existing order
// @Description Add new items to an existing order for current table session
// @Tags Order
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Param order body OrderItems true "Order items to add"
// @Success 201 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /current/items [post]
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

// GetOrderItems godoc
// @Summary Get order items for current session
// @Description Get all order items for the current table session with pagination
// @Tags Order
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Param page_number query int false "Page number for pagination" default(1)
// @Success 200 {object} middleware.SuccessResponse{data=domain.SearchCurrentOrderItemsResult}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /current/items [get]
func (s *Handler) GetOrderItems(c *fiber.Ctx) error {
	sessionID, err := getSession(c)
	if err != nil {
		return err
	}

	body := new(SearchCurrentOrderItems)
	if errValidate := c.QueryParser(body); errValidate != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, errValidate.Error())
	}

	if err = s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetCurrentOrderItems(c.Context(), sessionID, body.PageNumber)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order items success", result)
}

// GetOrderItemsByID godoc
// @Summary Get order item details by ID
// @Description Get specific order item details for current table session
// @Tags Order
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Param orderItemsID path string true "Order Item ID"
// @Success 200 {object} middleware.SuccessResponse{data=model.OrderItems}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /current/items/{orderItemsID} [get]
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

// UpdateOrderItemsStatusCancel godoc
// @Summary Cancel order item
// @Description Update order item status to cancelled for current table session
// @Tags Order
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Param orderItemsID path string true "Order Item ID"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /current/items/{orderItemsID}/status/cancel [patch]
func (s *Handler) UpdateOrderItemsStatusCancel(c *fiber.Ctx) error {
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

// SearchOrderItemsInComplete godoc
// @Summary Search incomplete order items
// @Description Search incomplete order items with filters
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search by name"
// @Param statusCode query []string false "Filter by status codes"
// @Param orderBy query string false "Order by field"
// @Param orderType query string false "Order direction (asc, desc)"
// @Success 200 {object} middleware.SuccessResponse{data=domain.SearchOrderItemsResult}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /{id}/items/status/incomplete [get]
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

// GetOrderItemsByTableID godoc
// @Summary Get order items for current session
// @Description Get all order items for the current table session with pagination
// @Tags Order
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Param page_number query int false "Page number for pagination" default(1)
// @Success 200 {object} middleware.SuccessResponse{data=domain.SearchCurrentOrderItemsResult}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /current/items [get]
func (s *Handler) GetOrderItemsByTableID(c *fiber.Ctx) error {
	tableID, err := utils.StrToInt64(c.Params("tableID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	body := new(SearchCurrentOrderItems)
	if errValidate := c.QueryParser(body); errValidate != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, errValidate.Error())
	}

	if err = s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	sessionID, customError := s.useCase.GetSessionIDByTableID(c.Context(), tableID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	result, customError := s.useCase.GetCurrentOrderItems(c.Context(), sessionID, body.PageNumber)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get order items success", result)
}

// UpdateOrderItemsStatusCancelByTableID godoc
// @Summary Cancel order item
// @Description Update order item status to cancelled for current table session
// @Tags Order
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Param orderItemsID path string true "Order Item ID"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /tables/{tableID}/items/{orderItemsID}/status/cancel [patch]
func (s *Handler) UpdateOrderItemsStatusCancelByTableID(c *fiber.Ctx) error {

	tableID, err := utils.StrToInt64(c.Params("tableID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	sessionID, customError := s.useCase.GetSessionIDByTableID(c.Context(), tableID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	orderItemsID, err := utils.StrToInt64(c.Params("orderItemsID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	customError = s.useCase.UpdateOrderItemsStatus(c.Context(), sessionID, shareModel.OrderItemsStatus{
		ID:         orderItemsID,
		StatusCode: "CANCELLED",
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "update order item status success", nil)
}

// UpdateOrderItemsStatusServedByTableID godoc
// @Summary Cancel order item
// @Description Update order item status to cancelled for current table session
// @Tags Order
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Param orderItemsID path string true "Order Item ID"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /tables/{tableID}/items/{orderItemsID}/status/cancel [patch]
func (s *Handler) UpdateOrderItemsStatusServedByTableID(c *fiber.Ctx) error {

	tableID, err := utils.StrToInt64(c.Params("tableID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	sessionID, customError := s.useCase.GetSessionIDByTableID(c.Context(), tableID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	orderItemsID, err := utils.StrToInt64(c.Params("orderItemsID"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	customError = s.useCase.UpdateOrderItemsStatus(c.Context(), sessionID, shareModel.OrderItemsStatus{
		ID:         orderItemsID,
		StatusCode: "SERVED",
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "update order item status success", nil)
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
