package http

import (
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"food-story/table-service/internal/domain"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func (s *Handler) GetTime(c *fiber.Ctx) error {
	result, customError := s.useCase.GetCurrentTime(c.Context())
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get current time success", result)
}

func (s *Handler) ListTableStatus(c *fiber.Ctx) error {
	result, customError := s.useCase.ListTableStatus(c.Context())
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get list task status success", result)
}

func (s *Handler) CreateTable(c *fiber.Ctx) error {
	body := new(Table)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.CreateTable(c.Context(), domain.Table{
		TableNumber: body.TableNumber,
		Seats:       body.Seats,
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseCreated(c, "create table success", createResponse{
		ID: strconv.FormatInt(result, 10),
	})
}

func (s *Handler) UpdateTable(c *fiber.Ctx) error {
	id, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	body := new(Table)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	customError := s.useCase.UpdateTable(c.Context(), domain.Table{
		ID:          id,
		TableNumber: body.TableNumber,
		Seats:       body.Seats,
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "update table success", nil)
}

func (s *Handler) UpdateTableStatus(c *fiber.Ctx) error {
	id, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	body := new(updateTableStatus)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	statusID, err := utils.StrToInt64(body.StatusID)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}
	customError := s.useCase.UpdateTableStatus(c.Context(), domain.TableStatus{
		ID:       id,
		StatusID: statusID,
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseCreated(c, "update table status success", nil)
}

func (s *Handler) SearchTable(c *fiber.Ctx) error {
	body := new(SearchTable)
	if err := c.QueryParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	payload := domain.SearchTables{
		NumberOfPeople: body.NumberOfPeople,
		StatusCode:     utils.FilterOutEmptyStr(body.Status),
		OrderBy:        body.OrderBy,
		OrderByType:    body.OrderByType,
		PageNumber:     body.PageNumber,
		PageSize:       body.PageSize,
	}

	if body.Search != "" {
		pareValue, err := strconv.ParseInt(body.Search, 10, 32)
		if err != nil {
			return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
		}
		value := int32(pareValue)
		payload.TableNumber = &value
	}

	if body.Seats != "" {
		pareValue, err := strconv.ParseInt(body.Seats, 10, 32)
		if err != nil {
			return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
		}
		value := int32(pareValue)
		payload.Seats = &value
	}

	result, customError := s.useCase.SearchTableByFilters(c.Context(), payload)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "search table success", result)
}

func (s *Handler) QuickSearchTable(c *fiber.Ctx) error {
	body := new(SearchTable)
	if err := c.QueryParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if body.NumberOfPeople == 0 {
		return middleware.ResponseError(fiber.StatusBadRequest, "number of people must be greater than 0")
	}

	payload := domain.SearchTables{
		NumberOfPeople: body.NumberOfPeople,
		OrderBy:        body.OrderBy,
		OrderByType:    body.OrderByType,
		PageNumber:     body.PageNumber,
		PageSize:       body.PageSize,
	}

	result, customError := s.useCase.QuickSearchAvailableTable(c.Context(), payload)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "search table success", result)
}

func (s *Handler) CreateTableSession(c *fiber.Ctx) error {
	body := new(TableSession)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	tableID, err := utils.StrToInt64(body.TableID)
	if err != nil {
		return middleware.ResponseError(fiber.StatusBadRequest, err.Error())
	}

	result, customError := s.useCase.CreateTableSession(c.Context(), domain.TableSession{
		TableID:        tableID,
		NumberOfPeople: body.NumberOfPeople,
	})
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseCreated(c, "create table success", createSessionResponse{
		URL: result,
	})
}

func (s *Handler) CurrentSession(c *fiber.Ctx) error {
	sessionIDData := c.Get("X-Session-Id")
	result, customError := s.useCase.GetCurrentSession(sessionIDData)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get current session success", result)
}
