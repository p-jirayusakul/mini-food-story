package http

import (
	"food-story/menu/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (s *Handler) CategoryList(c *fiber.Ctx) error {
	result, customError := s.useCase.ListCategory(c.Context())
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get list categories success", result)
}

func (s *Handler) CreateProduct(c *fiber.Ctx) error {
	body := new(Product)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	categoryID, err := utils.StrToInt64(body.CategoryID)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.CreateProduct(c.Context(), domain.Product{
		Name:        body.Name,
		NameEN:      body.NameEN,
		CategoryID:  categoryID,
		Price:       body.Price,
		Description: body.Description,
		IsAvailable: body.IsAvailable,
		ImageURL:    body.ImageURL,
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseCreated(c, "create table success", createResponse{
		ID: strconv.FormatInt(result, 10),
	})
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

//func (s *Handler) GetMenu(c *fiber.Ctx) error {
//	id, err := utils.StrToInt64(c.Params("id"))
//	if err != nil {
//		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
//	}
//
//
//}
