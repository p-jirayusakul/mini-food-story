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
}

func NewHTTPHandler(
	router fiber.Router,
	useCase usecase.PaymentUsecase,
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

	group.Post("/", s.CreatePaymentTransaction)
	group.Post("/callback", s.CallbackPaymentTransaction)

	group.Get("/methods", s.ListPaymentMethods)

}
