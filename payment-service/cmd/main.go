package main

import (
	"context"
	"fmt"
	"food-story/payment-service/internal"
	"food-story/pkg/common"
	"food-story/shared/config"
	"log"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"food-story/payment-service/docs"
)

func gracefulShutdown(fiberServer *internal.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// close all connection
	fiberServer.CloseAllConnection()

	if err := fiberServer.App.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	server := internal.New()
	port, _ := strconv.Atoi(server.Config.AppPort)
	initSwagger(server.Config)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		err := server.App.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}

func initSwagger(cfg config.Config) {
	port, _ := strconv.Atoi(cfg.AppPort)
	host := fmt.Sprintf("localhost:%d", port)

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Payment Service API"
	docs.SwaggerInfo.Description = "REST API for managing payment transactions and payment methods for restaurant orders"
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = cfg.BaseURL

	// dynamically set swagger info
	docs.SwaggerInfo.Host = host
	if strings.ToUpper(cfg.AppEnv) != common.DefaultAppEnv {
		docs.SwaggerInfo.Host = host // e.g. api.example.com
	}

	docs.SwaggerInfo.Schemes = []string{"http"}
	if strings.ToUpper(cfg.AppEnv) != common.DefaultAppEnv {
		docs.SwaggerInfo.Schemes = []string{"http"} // e.g. https
	}
}
