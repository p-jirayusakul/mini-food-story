package consumer

import (
	"context"
	"errors"
	"food-story/kitchen-service/internal/adapter/websocket"
	"github.com/IBM/sarama"
	"log"
)

type Consumer struct {
	Ready        chan bool
	WebSocketHub *websocket.Hub
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.Ready)
	return nil
}
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		// แจ้งแตือนข้อความ
		c.WebSocketHub.Broadcast <- msg.Value

		// แจ้งว่า message นี้ถูก consume แล้ว
		session.MarkMessage(msg, "")
	}

	return nil
}

func StartConsumer(ctx context.Context, topics []string, client sarama.ConsumerGroup, websocketHub *websocket.Hub) {
	consumer := Consumer{
		Ready:        make(chan bool),
		WebSocketHub: websocketHub,
	}

	go func() {
		for {
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()
	<-consumer.Ready
	log.Println("Sarama consumer up and running!...")
}
