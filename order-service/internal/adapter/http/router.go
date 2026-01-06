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

	// ใช้ middleware header ทีละ endpoint เพราะมีข้อจำกัดเรื่อง router group authentication
	groupCustomer := s.router.Group("/current", s.handleSessionID)
	groupCustomer.Post("/", s.CreateOrder)
	groupCustomer.Get("/", s.GetCurrentOrderByID)
	groupCustomer.Post("/items", s.CreateOrderItems)
	groupCustomer.Get("/items", s.GetCurrentOrderItems)
	groupCustomer.Get("/items/:orderItemsID<int>", s.GetCurrentOrderItemsByID)
	groupCustomer.Patch("/items/:orderItemsID<int>/status/cancel", s.UpdateCurrentOrderItemsStatusCancel)

	//groupStaff := s.router.Group("")

	s.router.Post("/", s.CreateOrderByStaff)
	s.router.Post("/:id<int>/items", s.CreateOrderItemsByStaff)
	s.router.Get("/:id<int>/items/status/incomplete", s.SearchOrderItemsInComplete)
	s.router.Get("/:id<int>/items", s.GetOrderItems)
	s.router.Patch("/:id<int>/items/:orderItemsID<int>/status/cancel", s.UpdateOrderItemsStatusCancel)
	s.router.Patch("/:id<int>/items/:orderItemsID<int>/status/serve", s.UpdateOrderItemsStatusServed)
}

func (s *Handler) handleSessionID(c *fiber.Ctx) error {
	sessionIDAny := c.Get("X-Session-Id")
	if sessionIDAny == "" {
		return middleware.ResponseError(c, exceptions.Error(exceptions.CodeUnauthorized, exceptions.ErrFailedToReadSession.Error()))
	}

	sessionID, err := utils.PareStringToUUID(sessionIDAny)
	if err != nil {
		return middleware.ResponseError(c, exceptions.Error(exceptions.CodeSystem, exceptions.ErrFailedToReadSession.Error()))
	}

	err = s.useCase.IsSessionValid(sessionID)
	if err != nil {
		return middleware.ResponseError(c, err)
	}

	c.Locals("sessionID", sessionID.String())
	return c.Next()
}
