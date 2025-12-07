package http

import (
	"food-story/menu-service/internal/usecase"
	"food-story/pkg/exceptions"
	"food-story/pkg/middleware"
	"food-story/pkg/utils"
	"food-story/shared/config"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	router    fiber.Router
	useCase   usecase.Usecase
	validator *middleware.CustomValidator
	config    config.Config
}

func NewHTTPHandler(
	router fiber.Router,
	useCase usecase.Usecase,
	validator *middleware.CustomValidator,
	config config.Config,

) *Handler {
	handler := &Handler{
		router,
		useCase,
		validator,
		config,
	}
	handler.setupRoutes()
	return handler
}

func (s *Handler) setupRoutes() {
	group := s.router.Group("/")

	groupCustomer := group.Group("/customer", s.handleSessionID)
	groupCustomer.Get("", s.SearchMenu)
	groupCustomer.Get("/:id<int>", s.GetProductByID)
	groupCustomer.Get("/category", s.CategoryList)
	groupCustomer.Get("/session/current", s.SessionCurrent)

	groupStaff := group.Group("/staff")
	groupStaff.Get("", s.SearchMenu)
	groupStaff.Get("/session-extension", s.ListProductTimeExtension)
	groupStaff.Get("/:id<int>", s.GetProductByID)
	groupStaff.Get("/category", s.CategoryList)
}

func (s *Handler) handleSessionID(c *fiber.Ctx) error {
	sessionIDAny := c.Get("X-Session-Id")
	if sessionIDAny == "" {
		return middleware.ResponseError(c, exceptions.Error(exceptions.CodeUnauthorized, exceptions.ErrFailedToReadSession.Error()))
	}

	sessionID, err := utils.PareStringToUUID(sessionIDAny)
	if err != nil {
		return middleware.ResponseError(c, exceptions.Error(exceptions.CodeSystem, exceptions.ErrFailedToReadSession.Error()))
	}

	err = s.useCase.IsSessionValid(sessionID)
	if err != nil {
		return middleware.ResponseError(c, err)
	}

	return c.Next()
}
