package producer

import (
	"encoding/json"
	"food-story/order-service/internal/domain"
	"food-story/shared/kafka"
	"github.com/IBM/sarama"
)

type QueueProducerInterface interface {
	PublishOrder(item domain.OrderItems) error
}
type OrderProducer struct {
	Producer sarama.SyncProducer
}

func NewQueue(producer sarama.SyncProducer) *OrderProducer {
	return &OrderProducer{
		producer,
	}
}

var _ QueueProducerInterface = (*OrderProducer)(nil)

func (p *OrderProducer) PublishOrder(item domain.OrderItems) error {
	message, err := json.Marshal(item)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: kafka.OrderItemsCreatedTopic,
		Value: sarama.ByteEncoder(message),
	}
	_, _, err = p.Producer.SendMessage(msg)
	return err
}
