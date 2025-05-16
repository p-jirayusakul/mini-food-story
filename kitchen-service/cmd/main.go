package main

import (
	"context"
	"fmt"
	"food-story/kitchen-service/internal"
	"food-story/kitchen-service/internal/adapter/queue/consumer"
	"food-story/shared/kafka"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

//func initKafka(ctx context.Context, websocketHub *websocket.Hub) sarama.ConsumerGroup {
//	topics := []string{kafka.OrderItemsCreatedTopic}
//	brokers := []string{"localhost:9092"}
//	client := kafka.InitConsumer("kitchen-group", brokers)
//	consumer.StartConsumer(ctx, topics, client, websocketHub)
//
//	return client
//}

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

	// close database
	fiberServer.CloseDB()
	log.Println("Database closed")

	// close consumer
	if err := fiberServer.KafkaClient.Close(); err != nil {
		log.Printf("close consumer error: %v", err)
	}
	cancelConsumer()
	log.Println("Kafka Consumer closed")

	// close websocket hub
	fiberServer.CloseWebsocketHub()
	log.Println("Websocket Hub closed")

	if err := fiberServer.App.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	server := internal.New()

	// start websocket hub
	go server.WebsocketHub.Run()

	// เริ่มต้น Kafka Consumer
	ctxConsumer, cancelConsumer := context.WithCancel(context.Background())
	go consumer.StartConsumer(ctxConsumer, []string{kafka.OrderItemsCreatedTopic}, server.KafkaClient, server.WebsocketHub)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
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
