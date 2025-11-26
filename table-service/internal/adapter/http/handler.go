package http

import (
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"food-story/table-service/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ListTableStatus godoc
// @Summary Get list of table status
// @Description Get list of all available table statuses
// @Tags Table
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} middleware.SuccessResponse{data=[]domain.Status}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /status [get]
func (s *Handler) ListTableStatus(c *fiber.Ctx) error {
	result, err := s.useCase.ListTableStatus(c.Context())
	if err != nil {
		return middleware.ResponseError(c, err)
	}

	return middleware.ResponseOK(c, result)
}

// ListSessionExtensionReason godoc
// @Summary Get list of table status
// @Description Get list of all available table statuses
// @Tags Table
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} middleware.SuccessResponse{data=[]domain.ListSessionExtensionReason}
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 401 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /status [get]
func (s *Handler) ListSessionExtensionReason(c *fiber.Ctx) error {
	result, err := s.useCase.ListSessionExtensionReason(c.Context())
	if err != nil {
		return middleware.ResponseError(c, err)
	}

	return middleware.ResponseOK(c, result)
}

// UpdateTableStatus godoc
// @Summary Update table status
// @Description Update status for existing table
// @Tags Table
// @Security BearerAuth
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
		return middleware.ResponseError(c, validateFail(err))
	}

	body := new(updateTableStatus)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}

	statusID, err := utils.StrToInt64(body.StatusID)
	if err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}
	err = s.useCase.UpdateTableStatus(c.Context(), domain.TableStatus{
		ID:       id,
		StatusID: statusID,
	})
	if err != nil {
		return middleware.ResponseError(c, err)
	}

	return middleware.ResponseOK(c, nil)
}

// SearchTable godoc
// @Summary Search table availability
// @Description Search tables by filters like number of people, table number, seats, and status
// @Tags Table
// @Security BearerAuth
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
		return middleware.ResponseError(c, validateFail(err))
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(c, validateFail(err))
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
			return middleware.ResponseError(c, validateFail(err))
		}
		value := int32(pareValue)
		payload.TableNumber = &value
	}

	if body.Seats != "" {
		pareValue, err := strconv.ParseInt(body.Seats, 10, 32)
		if err != nil {
			return middleware.ResponseError(c, validateFail(err))
		}
		value := int32(pareValue)
		payload.Seats = &value
	}

	result, err := s.useCase.SearchTableByFilters(c.Context(), payload)
	if err != nil {
		return middleware.ResponseError(c, err)
	}

	return middleware.ResponseOKWithPagination(c, middleware.ResponseWithPaginationPayload{
		PageNumber: result.PageNumber,
		PageSize:   result.PageSize,
		TotalItems: result.TotalItems,
		TotalPages: result.TotalPages,
		Data:       result.Data,
	})
}

// QuickSearchTable godoc
// @Summary Quick search for available tables
// @Description Quickly search for available tables based on number of people
// @Tags Table
// @Security BearerAuth
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
		return middleware.ResponseError(c, validateFail(err))
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}

	payload := domain.SearchTables{
		NumberOfPeople: body.NumberOfPeople,
		OrderBy:        body.OrderBy,
		OrderByType:    body.OrderByType,
		PageNumber:     body.PageNumber,
		PageSize:       body.PageSize,
	}

	result, err := s.useCase.QuickSearchAvailableTable(c.Context(), payload)
	if err != nil {
		return middleware.ResponseError(c, err)
	}

	return middleware.ResponseOKWithPagination(c, middleware.ResponseWithPaginationPayload{
		PageNumber: result.PageNumber,
		PageSize:   result.PageSize,
		TotalItems: result.TotalItems,
		TotalPages: result.TotalPages,
		Data:       result.Data,
	})
}

// CreateTableSession godoc
// @Summary Create new table session
// @Description Create a new session for a table with specified number of people
// @Tags Table
// @Security BearerAuth
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
		return middleware.ResponseError(c, validateFail(err))
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}

	tableID, err := utils.StrToInt64(body.TableID)
	if err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}

	result, err := s.useCase.CreateTableSession(c.Context(), domain.TableSession{
		TableID:        tableID,
		NumberOfPeople: body.NumberOfPeople,
	})
	if err != nil {
		return middleware.ResponseError(c, err)
	}

	return middleware.ResponseCreated(c, createSessionResponse{
		URL: result,
	})
}

// UpdateTableStatusAvailable godoc
// @Summary Update table status
// @Description Update status for existing table
// @Tags Table
// @Security BearerAuth
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
// @Router /{id}/status/available [patch]
func (s *Handler) UpdateTableStatusAvailable(c *fiber.Ctx) error {
	id, err := utils.StrToInt64(c.Params("id"))
	if err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}

	err = s.useCase.UpdateTableStatusAvailable(c.Context(), id)
	if err != nil {
		return middleware.ResponseError(c, err)
	}

	return middleware.ResponseOK(c, nil)
}

// SessionExtension godoc
// @Summary Update table status
// @Description Update status for existing table
// @Tags Table
// @Security BearerAuth
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
// @Router /{id}/status/available [patch]
func (s *Handler) SessionExtension(c *fiber.Ctx) error {

	body := new(SessionExtensionRequest)
	if err := c.BodyParser(body); err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}

	if err := s.validator.Validate(body); err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}

	tableID, err := utils.StrToInt64(body.TableID)
	if err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}

	productID, err := utils.StrToInt64(body.ProductID)
	if err != nil {
		return middleware.ResponseError(c, validateFail(err))
	}

	err = s.useCase.SessionExtension(c.Context(), domain.SessionExtension{
		TableID:    tableID,
		ProductID:  productID,
		ReasonCode: body.ReasonCode,
	})
	if err != nil {
		return middleware.ResponseError(c, err)
	}
	return middleware.ResponseOK(c, nil)
}

func validateFail(err error) error {
	return exceptions.Error(exceptions.CodeBusiness, err.Error())
}
