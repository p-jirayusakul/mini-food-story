package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"food-story/kitchen-service/internal/domain"
	"github.com/IBM/sarama"
	"log"
)

type Consumer struct {
	Ready chan bool
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.Ready)
	return nil
}
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		var orderItems domain.OrderItems
		if err := json.Unmarshal(msg.Value, &orderItems); err != nil {
			continue
		}

		// ประมวลผลคำสั่งอาหาร
		processOrder(orderItems)

		// แจ้งว่า message นี้ถูก consume แล้ว
		session.MarkMessage(msg, "")
	}

	return nil
}

func StartConsumer(ctx context.Context, topics []string, client sarama.ConsumerGroup) {
	consumer := Consumer{
		Ready: make(chan bool),
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

func processOrder(order domain.OrderItems) {
	log.Printf("  - Menu: %d %s x %d\n", order.OrderID, order.ProductName, order.Quantity)
}
