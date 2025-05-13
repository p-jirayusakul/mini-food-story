package http

import (
	"food-story/order/internal/usecase"
	"food-story/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	router    fiber.Router
	useCase   usecase.Usecase
	validator *middleware.CustomValidator
}

func NewHTTPHandler(
	router fiber.Router,
	useCase usecase.Usecase,
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
	group := s.router.Group("/orders")

	group.Post("", s.CreateOrder)
	group.Get("/:id<int>", s.GetOrderByID)
	group.Patch("/:id<int>/status", s.UpdateOrderStatus)

	group.Post("/:id<int>/items", s.CreateOrderItems)
	group.Get("/:id<int>/items", s.GetOrderItems)
	group.Get("/:id<int>/items/:orderItemsID<int>", s.GetOrderItemsByID)
	group.Patch("/:id<int>/items/:orderItemsID<int>/status", s.UpdateOrderItemsStatus)

}
