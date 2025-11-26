package http

import (
	"food-story/pkg/middleware"
	"food-story/table-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	router    fiber.Router
	useCase   usecase.UseCase
	validator *middleware.CustomValidator
	auth      middleware.AuthInterface
}

func NewHTTPHandler(
	router fiber.Router,
	useCase usecase.UseCase,
	validator *middleware.CustomValidator,
	auth middleware.AuthInterface,
) *Handler {
	handler := &Handler{
		router,
		useCase,
		validator,
		auth,
	}
	handler.setupRoutes()
	return handler
}

func (s *Handler) setupRoutes() {

	group := s.router.Group("/")
	group.Use(s.auth.JWTMiddleware(), s.auth.RequireRole([]string{"CASHIER", "WAITER"}))

	group.Get("", s.SearchTable)
	group.Get("/quick-search", s.QuickSearchTable)

	group.Get("/status", s.ListTableStatus)
	group.Patch("/:id<int>/status/available", s.UpdateTableStatusAvailable)

	group.Post("/session", s.CreateTableSession)
	group.Post("/session/extension", s.SessionExtension)
	group.Get("/session/reason", s.ListSessionExtensionReason)

}
