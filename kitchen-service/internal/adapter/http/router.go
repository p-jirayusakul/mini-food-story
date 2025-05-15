package http

import (
	"food-story/kitchen-service/internal/usecase"
	"food-story/pkg/middleware"
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
	group := s.router.Group("/orders")
	group.Patch("/:id<int>/items/:orderItemsID<int>/status", s.UpdateOrderItemsStatus)

}
