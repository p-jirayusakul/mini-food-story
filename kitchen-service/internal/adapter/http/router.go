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
	group := s.router.Group("/orders")
	group.Use(s.auth.JWTMiddleware())

	authGroup := group.Group("", s.auth.RequireRole([]string{"KITCHEN"}))
	authGroup.Get("/search/items", s.SearchOrderItems)
	authGroup.Get("/:id<int>/items", s.GetOrderItems)
	authGroup.Get("/:id<int>/items/:orderItemsID<int>", s.GetOrderItemsByID)

	authGroup.Patch("/:id<int>/items/:orderItemsID<int>/status/serve", s.UpdateOrderItemsStatusServe)
	authGroup.Patch("/:id<int>/items/:orderItemsID<int>/status/cancel", s.UpdateOrderItemsStatusCancel)

}
