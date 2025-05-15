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
	group.Get("/kitchen", websocket.New(func(c *websocket.Conn) {
		s.hub.Register <- c

		defer func() {
			s.hub.Unregister <- c
		}()

		// ไม่ต้องอ่าน message จาก client ถ้าไม่จำเป็น
		for {
			// แค่ wait ให้ client disconnect ไปเอง
			if _, _, err := c.ReadMessage(); err != nil {
				break
			}
		}
	}))
}

func (s *WSHandler) WS(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
