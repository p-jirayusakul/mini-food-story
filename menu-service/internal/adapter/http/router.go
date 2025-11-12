package http

import (
	"fmt"
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
	auth      middleware.AuthInterface
}

func NewHTTPHandler(
	router fiber.Router,
	useCase usecase.Usecase,
	validator *middleware.CustomValidator,
	config config.Config,
	auth middleware.AuthInterface,

) *Handler {
	handler := &Handler{
		router,
		useCase,
		validator,
		config,
		auth,
	}
	handler.setupRoutes()
	return handler
}

func (s *Handler) setupRoutes() {
	group := s.router.Group("/")

	group.Get("", middleware.CheckSessionHeader(s.config.SecretKey), s.handleSessionID, s.SearchMenu)
	group.Get("/:id<int>", middleware.CheckSessionHeader(s.config.SecretKey), s.handleSessionID, s.GetProductByID)
	group.Get("/category", middleware.CheckSessionHeader(s.config.SecretKey), s.handleSessionID, s.CategoryList)
	group.Get("/session/current", middleware.CheckSessionHeader(s.config.SecretKey), s.handleSessionID, s.SessionCurrent)

	roles := []string{"CASHIER", "WAITER"}
	group.Get("/internal", s.auth.JWTMiddleware(), s.auth.RequireRole(roles), s.SearchMenu)
	group.Get("/internal/:id<int>", s.auth.JWTMiddleware(), s.auth.RequireRole(roles), s.GetProductByID)
	group.Get("/internal/category", s.auth.JWTMiddleware(), s.auth.RequireRole(roles), s.CategoryList)
}

func (s *Handler) handleSessionID(c *fiber.Ctx) error {
	sessionIDAny, ok := c.Locals("sessionID").(string)
	if !ok {
		fmt.Println("1")
		return middleware.ResponseError(fiber.StatusInternalServerError, exceptions.ErrFailedToReadSession.Error())
	}

	sessionID, err := utils.PareStringToUUID(sessionIDAny)
	if err != nil {
		fmt.Println("2")

		return middleware.ResponseError(fiber.StatusInternalServerError, exceptions.ErrFailedToReadSession.Error())
	}

	customError := s.useCase.IsSessionValid(sessionID)
	if customError != nil {
		fmt.Println("3")

		return middleware.ResponseError(exceptions.MapToHTTPStatusCode(customError.Status), customError.Errors.Error())
	}

	return c.Next()
}
