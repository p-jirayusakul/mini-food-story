package websocket

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (s *WSHandler) WS(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (s *WSHandler) WSKitchen(c *websocket.Conn) {
	s.hub.Register <- c
	defer func() {
		s.hub.Unregister <- c
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break // client ปิดเอง
		}
	}
}
