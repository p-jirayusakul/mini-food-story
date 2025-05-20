package http

import (
	"food-story/order-service/internal/usecase"
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

	withoutSession := s.router.Group("/orders")
	withoutSession.Get("/:id<int>/items/status/incomplete", s.SearchOrderItemsInComplete)

	group := s.router.Group("/orders", middleware.CheckSessionHeader(s.config.SecretKey), s.handleSessionID)
	group.Post("/current", s.CreateOrder)
	group.Get("/current", s.GetOrderByID)
	group.Post("/current/items", s.CreateOrderItems)
	group.Get("/current/items", s.GetOrderItems)
	group.Get("/current/items/:orderItemsID<int>", s.GetOrderItemsByID)
	group.Patch("/current/items/:orderItemsID<int>/status/cancelled", s.UpdateOrderItemsStatusCancelled)

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
