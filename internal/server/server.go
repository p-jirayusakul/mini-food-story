package server

import (
	"github.com/gofiber/fiber/v2"

	"food-story/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "food-story",
			AppName:      "food-story",
		}),

		db: database.New(),
	}

	return server
}
