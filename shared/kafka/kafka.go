package kafka

import (
	"github.com/IBM/sarama"
	"log"
	"time"
)

func InitConsumer(group string, brokers []string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = sarama.V3_9_0_0 // confluentinc 7.9.X
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = true

	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	return client
}

func InitProducer(brokers []string) (producer sarama.SyncProducer, client sarama.Client, err error) {

	config := sarama.NewConfig()
	config.Version = sarama.V3_9_0_0 // confluentinc 7.9.X
	config.Net.DialTimeout = 5 * time.Second
	config.Producer.Return.Successes = true

	// สร้าง client
	client, err = sarama.NewClient(brokers, config)
	if err != nil {
		return nil, nil, err
	}

	// สร้าง producer จาก client
	producer, err = sarama.NewSyncProducerFromClient(client)
	if err != nil {
		_ = client.Close()
		return nil, nil, err
	}

	return producer, client, nil
}
