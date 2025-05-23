package http

import (
	"food-story/menu-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// CategoryList godoc
// @Summary Get list of categories
// @Description Get list of all available product categories
// @Tags Category
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Success 200 {object} middleware.SuccessResponse{data=[]domain.Category}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /category [get]
func (s *Handler) CategoryList(c *fiber.Ctx) error {

	result, customError := s.useCase.ListCategory(c.Context())
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get list category success", result)
}

// SearchMenu godoc
// @Summary Search menu items
// @Description Search menu items with filters
// @Tags Menu
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search by name"
// @Param categoryID query []string false "Filter by category IDs"
// @Param orderBy query string false "Order by field (id, tableNumber, seats, status)"
// @Param orderType query string false "Order direction (asc, desc)"
// @Success 200 {object} middleware.SuccessResponse{data=domain.SearchProductResult}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router / [get]
func (s *Handler) SearchMenu(c *fiber.Ctx) error {
	body := new(SearchMenu)
	if err := c.QueryParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	payload := domain.SearchProduct{
		Name:        body.Search,
		CategoryID:  utils.FilterOutZero(body.CategoryID),
		IsAvailable: true,
		OrderByType: body.OrderByType,
		OrderBy:     body.OrderBy,
		PageSize:    body.PageSize,
		PageNumber:  body.PageNumber,
	}

	result, customError := s.useCase.SearchProductByFilters(c.Context(), payload)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "search menu success", result)
}

// GetProductByID godoc
// @Summary Get menu item by ID
// @Description Get menu item details by product ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Param id path string true "Product ID"
// @Success 200 {object} middleware.SuccessResponse{data=domain.Product}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /{id} [get]
func (s *Handler) GetProductByID(c *fiber.Ctx) error {
	productID, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetProductByID(c.Context(), productID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get menu success", result)
}
