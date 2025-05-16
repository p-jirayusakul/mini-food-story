package websocket

import (
	"food-story/shared/config"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WSHandler struct {
	router fiber.Router
	config config.Config
	hub    *Hub
}

func NewWSHandler(
	router fiber.Router,
	config config.Config,
	hub *Hub,
) *WSHandler {
	handler := &WSHandler{
		router,
		config,
		hub,
	}
	handler.setupRoutes()
	return handler
}

func (s *WSHandler) setupRoutes() {
	group := s.router.Group("/ws")
	group.Use("/ws", s.WS)
	group.Get("/kitchen", websocket.New(s.WSKitchen))
}
