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
	group := s.router.Group("/")
	group.Post("/callback", s.CallbackPaymentTransaction)

	group.Post("/", s.CreatePaymentTransaction)
	group.Get("/methods", s.ListPaymentMethods)
	group.Get("/transactions/:transactionID/stream", s.StreamPaymentStatusByTransaction)
	group.Get("/transactions/:transactionID/qr", s.PaymentTransactionQR)
}
