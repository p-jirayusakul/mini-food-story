package http

import (
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"food-story/table-service/internal/domain"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// ListTableStatus godoc
// @Summary Get list of table status
// @Description Get list of all available table statuses
// @Tags Table
// @Accept json
// @Produce json
// @Success 200 {object} middleware.SuccessResponse{data=[]domain.Status}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /status [get]
func (s *Handler) ListTableStatus(c *fiber.Ctx) error {
	result, customError := s.useCase.ListTableStatus(c.Context())
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get list task status success", result)
}

// CreateTable godoc
// @Summary Create new table
// @Description Create a new table with specified number and seats
// @Tags Table
// @Accept json
// @Produce json
// @Param table body Table true "Table details"
// @Success 201 {object} middleware.SuccessResponse{data=createResponse}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router / [post]
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

// UpdateTable godoc
// @Summary Update table details
// @Description Update table number and seats for existing table
// @Tags Table
// @Accept json
// @Produce json
// @Param id path string true "Table ID"
// @Param table body Table true "Table details"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /{id} [put]
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

// UpdateTableStatus godoc
// @Summary Update table status
// @Description Update status for existing table
// @Tags Table
// @Accept json
// @Produce json
// @Param id path string true "Table ID"
// @Param status body updateTableStatus true "Table status details"
// @Success 200 {object} middleware.SuccessResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /{id}/status [patch]
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

	return middleware.ResponseOK(c, "update table status success", nil)
}

// SearchTable godoc
// @Summary Search table availability
// @Description Search tables by filters like number of people, table number, seats, and status
// @Tags Table
// @Accept json
// @Produce json
// @Param numberOfPeople query int false "Number of people"
// @Param search query string false "Search by table number"
// @Param seats query string false "Filter by seats"
// @Param status query []string false "Filter by status codes"
// @Param pageNumber query int false "Page number for pagination"
// @Param pageSize query int false "Page size for pagination"
// @Param orderBy query string false "Order by field (id, tableNumber, seats, status)"
// @Param orderByType query string false "Order direction (asc, desc)"
// @Success 200 {object} middleware.SuccessResponse{data=domain.SearchTablesResult}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router / [get]
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

// QuickSearchTable godoc
// @Summary Quick search for available tables
// @Description Quickly search for available tables based on number of people
// @Tags Table
// @Accept json
// @Produce json
// @Param numberOfPeople query int true "Number of people required"
// @Param pageNumber query int false "Page number for pagination"
// @Param pageSize query int false "Page size for pagination"
// @Param orderBy query string false "Order by field (id, tableNumber, seats, status)"
// @Param orderByType query string false "Order direction (asc, desc)"
// @Success 200 {object} middleware.SuccessResponse{data=domain.SearchTablesResult}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /quick-search [get]
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

// CreateTableSession godoc
// @Summary Create new table session
// @Description Create a new session for a table with specified number of people
// @Tags Table
// @Accept json
// @Produce json
// @Param table body TableSession true "Table session details"
// @Success 201 {object} middleware.SuccessResponse{data=createSessionResponse}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /session [post]
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

// CurrentSession godoc
// @Summary Get current table session
// @Description Get details of the current active table session
// @Tags Table
// @Accept json
// @Produce json
// @Param X-Session-Id header string true "Session ID"
// @Success 200 {object} middleware.SuccessResponse{data=model.CurrentTableSession}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /session/current [get]
func (s *Handler) CurrentSession(c *fiber.Ctx) error {
	sessionIDData := c.Get("X-Session-Id")
	result, customError := s.useCase.GetCurrentSession(sessionIDData)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return middleware.ResponseOK(c, "get current session success", result)
}
