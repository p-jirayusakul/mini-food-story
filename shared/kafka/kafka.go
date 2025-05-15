package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

func InitConsumer(group string, brokers []string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = sarama.V3_9_0_0 // confluentinc 7.9.X
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	return client
}
