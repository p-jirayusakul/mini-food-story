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
	group.Get("/search/items", s.SearchOrderItems)
	group.Get("/:id<int>/items", s.GetOrderItems)
	group.Get("/:id<int>/items/:orderItemsID<int>", s.GetOrderItemsByID)

	group.Patch("/:id<int>/items/:orderItemsID<int>/status/serve", s.UpdateOrderItemsStatusServe)
	group.Patch("/:id<int>/items/:orderItemsID<int>/status/cancel", s.UpdateOrderItemsStatusCancel)

}
