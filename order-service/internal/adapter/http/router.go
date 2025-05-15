package http

import (
	"food-story/order-service/internal/usecase"
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

	group.Post("/current", s.CreateOrder)
	group.Get("/current", s.GetOrderByID)

	group.Post("/current/items", s.CreateOrderItems)
	group.Get("/current/items", s.GetOrderItems)
	group.Get("/current/items/:orderItemsID<int>", s.GetOrderItemsByID)
	group.Patch("/current/items/:orderItemsID<int>/status/cancelled", s.UpdateOrderItemsStatusCancelled)

}
