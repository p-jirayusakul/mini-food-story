package http

import (
	"food-story/payment-service/internal/usecase"
	"food-story/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	router    fiber.Router
	useCase   usecase.PaymentUsecase
	validator *middleware.CustomValidator
	auth      middleware.AuthInterface
}

func NewHTTPHandler(
	router fiber.Router,
	useCase usecase.PaymentUsecase,
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
	const role = "CASHIER"
	group := s.router.Group("/")
	group.Post("/callback", s.CallbackPaymentTransaction)

	group.Post("/", s.auth.JWTMiddleware(), s.auth.RequireRole([]string{role}), s.CreatePaymentTransaction)
	group.Get("/methods", s.auth.JWTMiddleware(), s.auth.RequireRole([]string{role}), s.ListPaymentMethods)
}
