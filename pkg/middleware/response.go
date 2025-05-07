package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type SuccessResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"do something completed"`
	Data    interface{} `json:"data" swaggertype:"object"`
}

type ErrorResponse struct {
	Status  string      `json:"status" example:"error"`
	Message string      `json:"message" example:"something went wrong"`
	Data    interface{} `json:"data" swaggertype:"object"`
}

func ResponseOK(c *fiber.Ctx, message string, payload interface{}) error {
	return c.Status(fiber.StatusOK).JSON(SuccessResponse{Message: message, Status: "success", Data: payload})
}

func ResponseCreated(c *fiber.Ctx, message string, payload interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(SuccessResponse{Message: message, Status: "success", Data: payload})
}

func ResponseError(httpCode int, message string) error {
	if httpCode >= 500 {
		slog.Error(message)
	}
	return fiber.NewError(httpCode, message)
}
