package middleware

import "github.com/gofiber/fiber/v2/middleware/cors"

func DefaultCorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Connection",
		AllowMethods: "GET, PUT, POST, PATCH, DELETE, OPTIONS",
	}
}
