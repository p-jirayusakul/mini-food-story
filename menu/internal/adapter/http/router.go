package http

import (
	"food-story/menu/internal/usecase"
	"food-story/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	router    fiber.Router
	useCase   usecase.ProductUsecase
	validator *middleware.CustomValidator
}

func NewHTTPHandler(
	router fiber.Router,
	useCase usecase.ProductUsecase,
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
	group := s.router.Group("/menu")

	group.Get("", s.SearchMenu)
	//group.Put("/:id<int>", s.GetMenu)
	group.Get("/category", s.CategoryList)

	groupProduct := s.router.Group("/product")
	groupProduct.Post("", s.CreateProduct)

}
