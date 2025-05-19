package middleware

import (
	"food-story/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func CheckSessionHeader(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionIDData := c.Get("X-Session-Id")
		if sessionIDData == "" {
			return ResponseError(fiber.StatusUnauthorized, "X-Session-Id header is required")
		}

		sessionID, err := utils.DecryptSessionToUUID(sessionIDData, []byte(secretKey))
		if err != nil {
			return ResponseError(fiber.StatusForbidden, "Invalid session")
		}

		c.Locals("sessionID", sessionID.String())
		return c.Next()
	}
}
