package main

import (
	"context"
	"fmt"
	"food-story/kitchen-service/internal"
	"food-story/kitchen-service/internal/adapter/queue/consumer"
	"food-story/pkg/common"
	"food-story/shared/config"
	"food-story/shared/kafka"
	"log"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"food-story/kitchen-service/docs"
)

func gracefulShutdown(fiberServer *internal.FiberServer, cancelConsumer context.CancelFunc, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has minute 5 to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// close consumer
	if err := fiberServer.KafkaConsumer.Close(); err != nil {
		log.Printf("close consumer error: %v", err)
	}
	cancelConsumer()
	log.Println("Kafka Consumer closed")

	// close database
	fiberServer.CloseAllConnection()

	if err := fiberServer.App.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	server := internal.New()
	port, _ := strconv.Atoi(server.Config.AppPort)
	initSwagger(server.Config)

	// start websocket hub
	go server.WebsocketHub.Run()

	// เริ่มต้น Kafka Consumer
	ctxConsumer, cancelConsumer := context.WithCancel(context.Background())
	go consumer.Run(ctxConsumer, []string{kafka.OrderItemsCreatedTopic}, server.KafkaConsumer, server.WebsocketHub)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		err := server.App.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, cancelConsumer, done)

	<-done
	log.Println("Graceful shutdown complete.")
}

func initSwagger(cfg config.Config) {
	port, _ := strconv.Atoi(cfg.AppPort)
	host := fmt.Sprintf("localhost:%d", port)

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Kitchen Service API"
	docs.SwaggerInfo.Description = "REST API for managing restaurant kitchen operations, including order item processing, status updates and kitchen notifications"
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
