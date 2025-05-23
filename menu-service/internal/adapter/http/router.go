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
	group := s.router.Group("/", middleware.CheckSessionHeader(s.config.SecretKey), s.handleSessionID)

	group.Get("", s.SearchMenu)
	group.Get("/:id<int>", s.GetProductByID)
	group.Get("/category", s.CategoryList)
}

func (s *Handler) handleSessionID(c *fiber.Ctx) error {
	sessionIDAny, ok := c.Locals("sessionID").(string)
	if !ok {
		return middleware.ResponseError(fiber.StatusInternalServerError, exceptions.ErrFailedToReadSession.Error())
	}

	sessionID, err := utils.PareStringToUUID(sessionIDAny)
	if err != nil {
		return middleware.ResponseError(fiber.StatusInternalServerError, exceptions.ErrFailedToReadSession.Error())
	}

	customError := s.useCase.IsSessionValid(sessionID)
	if customError != nil {
		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return c.Next()
}
