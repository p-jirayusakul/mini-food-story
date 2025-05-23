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
}

func NewHTTPHandler(
	router fiber.Router,
	useCase usecase.UseCase,
	validator *middleware.CustomValidator,
) *Handler {
	handler := &Handler{
		router,
		useCase,
		validator,
	}
	handler.setupRoutes()
	return handler
}

func (s *Handler) setupRoutes() {
	group := s.router.Group("/")

	group.Get("/status", s.ListTableStatus)
	group.Post("", s.CreateTable)
	group.Post("/session", s.CreateTableSession)
	group.Get("/session/current", s.CurrentSession)

	group.Put("/:id<int>", s.UpdateTable)
	group.Patch("/:id<int>/status", s.UpdateTableStatus)
	group.Get("", s.SearchTable)
	group.Get("/quick-search", s.QuickSearchTable)

}
