package middleware

import (
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func CheckSessionHeader(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionIDData := c.Get("X-Session-Id")
		if sessionIDData == "" {
			return ResponseError(c, exceptions.Error(exceptions.CodeUnauthorized, "X-Session-Id header is required"))
		}

		sessionID, err := utils.DecryptSessionToUUID(sessionIDData, []byte(secretKey))
		if err != nil {
			return ResponseError(c, exceptions.Error(exceptions.CodeUnauthorized, "invalid session"))
		}

		c.Locals("sessionID", sessionID.String())
		return c.Next()
	}
}
