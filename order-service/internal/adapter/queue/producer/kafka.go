package producer

import (
	"encoding/json"
	"food-story/shared/kafka"
	shareModel "food-story/shared/model"

	"github.com/IBM/sarama"
)

type QueueProducerInterface interface {
	PublishOrder(item shareModel.OrderItems) error
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

func (p *OrderProducer) PublishOrder(item shareModel.OrderItems) error {
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
