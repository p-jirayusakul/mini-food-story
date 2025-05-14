package producer

import (
	"encoding/json"
	"food-story/order-service/internal/domain"
	"github.com/IBM/sarama"
)

type OrderProducer struct {
	Producer sarama.SyncProducer
}

func NewOrderProducer(brokers []string) (*OrderProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &OrderProducer{Producer: producer}, nil
}

func (p *OrderProducer) PublishOrder(item domain.OrderItems) error {

	message, err := json.Marshal(item)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "order.items.placed",
		Value: sarama.ByteEncoder(message),
	}
	_, _, err = p.Producer.SendMessage(msg)
	return err
}
