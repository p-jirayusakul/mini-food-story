package middleware

import (
	"food-story/pkg/common"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LogHandler() func(*fiber.Ctx) error {
	return logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			return (c.Path() == common.BasePath+common.LivenessEndpoint) || (c.Path() == common.BasePath+common.ReadinessEndpoint)
		},
		Format:     "${time} | ${latency} | ${ip}:${port} -  ${status} ${method} ${path} | ${error}\n",
		TimeFormat: "2006/01/02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	})
}
