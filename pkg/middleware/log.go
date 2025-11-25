package middleware

import (
	"food-story/pkg/common"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LogHandler(basePath string) func(*fiber.Ctx) error {
	timeZone := "Asia/Bangkok"
	if os.Getenv("TZ") != "" {
		timeZone = os.Getenv("TZ")
	}

	return logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			path := c.Path()
			switch {
			case path == basePath+common.LivenessEndpoint:
				return true
			case path == basePath+common.ReadinessEndpoint:
				return true
			case strings.Contains(path, basePath+common.SwaggerEndpoint):
				return true
			default:
				return false
			}
		},
		Format:     "${time} | ${locals:requestid} | INBOUND | ${latency} | ${ip}:${port} - ${status} ${method} ${path} | ${error}\n",
		TimeFormat: "2006/01/02 15:04:05",
		TimeZone:   timeZone,
	})
}
