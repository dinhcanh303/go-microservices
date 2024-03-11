package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/wire"
	"golang.org/x/exp/slog"
)

const (
	_autoOffsetReset  = "earliest"
	_groupId          = "group"
	_bootstrapServers = "localhost"
)

type kafkaConsumer struct {
	consumer *kafka.Consumer
}

var _ KafkaConsumer = (*kafkaConsumer)(nil)

var KafkaConsumerSet = wire.NewSet(NewKafkaConsumer)

func NewKafkaConsumer() (KafkaConsumer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers": _bootstrapServers,
		"group.id":          _groupId,
		"auto.offset.reset": _autoOffsetReset,
	}
	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, err
	}
	return &kafkaConsumer{
		consumer: consumer,
	}, nil
}

// Close implements KafkaConsumer.
func (c *kafkaConsumer) Close() error {
	err := c.consumer.Close()
	if err != nil {
		return err
	}
	return nil
}

// ReadMessage implements KafkaConsumer.
func (c *kafkaConsumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	msg, err := c.consumer.ReadMessage(timeout)
	if err != nil {
		slog.Error("Consumer error: %v (%v)\n", err, msg)
		return nil, err
	}
	return msg, err
}

// Subscribe implements KafkaConsumer.
func (c *kafkaConsumer) Subscribe(topic string, rebalanceCb kafka.RebalanceCb) error {
	err := c.consumer.Subscribe(topic, rebalanceCb)
	if err != nil {
		return err
	}
	return nil
}

// SubscribeTopics implements KafkaConsumer.
func (c *kafkaConsumer) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error {
	err := c.consumer.SubscribeTopics(topics, rebalanceCb)
	if err != nil {
		return err
	}
	return nil
}
