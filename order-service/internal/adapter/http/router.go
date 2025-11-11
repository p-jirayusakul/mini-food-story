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
	secretKey := s.config.SecretKey

	// ใช้ middleware header ทีละ endpoint เพราะมีข้อจำกัดเรื่อง router group authentication
	s.router.Post("/current", middleware.CheckSessionHeader(secretKey), s.handleSessionID, s.CreateOrder)
	s.router.Get("/current", middleware.CheckSessionHeader(secretKey), s.handleSessionID, s.GetCurrentOrderByID)
	s.router.Post("/current/items", middleware.CheckSessionHeader(secretKey), s.handleSessionID, s.CreateOrderItems)
	s.router.Get("/current/items", middleware.CheckSessionHeader(secretKey), s.handleSessionID, s.GetCurrentOrderItems)
	s.router.Get("/current/items/:orderItemsID<int>", middleware.CheckSessionHeader(secretKey), s.handleSessionID, s.GetCurrentOrderItemsByID)
	s.router.Patch("/current/items/:orderItemsID<int>/status/cancel", middleware.CheckSessionHeader(secretKey), s.handleSessionID, s.UpdateCurrentOrderItemsStatusCancel)

	var roles = []string{"WAITER", "CASHIER"}
	s.router.Post("/:id<int>", s.auth.JWTMiddleware(), s.auth.RequireRole(roles), s.CreateOrderByStaff)
	s.router.Post("/:id<int>/items", s.auth.JWTMiddleware(), s.auth.RequireRole(roles), s.CreateOrderItemsByStaff)
	s.router.Get("/:id<int>/items/status/incomplete", s.auth.JWTMiddleware(), s.auth.RequireRole(roles), s.SearchOrderItemsInComplete)
	s.router.Get("/:id<int>/items", s.auth.JWTMiddleware(), s.auth.RequireRole(roles), s.GetOrderItems)
	s.router.Patch("/:id<int>/items/:orderItemsID<int>/status/cancel", s.auth.JWTMiddleware(), s.auth.RequireRole(roles), s.UpdateOrderItemsStatusCancel)
	s.router.Patch("/:id<int>/items/:orderItemsID<int>/status/serve", s.auth.JWTMiddleware(), s.auth.RequireRole(roles), s.UpdateOrderItemsStatusServed)
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
