package http

import (
	"food-story/menu-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func (s *Handler) CategoryList(c *fiber.Ctx) error {
	result, customError := s.useCase.ListCategory(c.Context())
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get list category success", result)
}

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

func (s *Handler) GetMenuByID(c *fiber.Ctx) error {
	id, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.GetProductByID(c.Context(), id)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get menu success", result)
}
