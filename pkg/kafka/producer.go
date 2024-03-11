package kafka

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/wire"
)

// RequiredAcks enum
const (
	// Define Maximum Of Retry When Write Message
	_defaultRetries = 5
	// Define Time Sleep When Write Fail
	_defaultKafkaSleep = 250 * time.Millisecond
)

type kafkaProducer struct {
	producer *kafka.Producer
}

var _ KafkaProducer = (*kafkaProducer)(nil)

var KafkaProducerSet = wire.NewSet(NewKafkaProducer)

func NewKafkaProducer() (*kafkaProducer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers": _bootstrapServers,
		"group.id":          _groupId,
		"auto.offset.reset": _autoOffsetReset,
	}
	producer, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{
		producer: producer,
	}, nil
}

// Produce implements KafkaProducer.
func (p *kafkaProducer) Produce(schema *Schema[Payload]) error {
	jsonData, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	key, err := json.Marshal(schema.Key)
	if err != nil {
		return err
	}
	err = p.producer.Produce(&kafka.Message{
		Key:       key,
		Value:     jsonData,
		Timestamp: schema.Timestamp,
		TopicPartition: kafka.TopicPartition{
			Partition: kafka.PartitionAny,
		},
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Close implements KafkaConsumer.
func (p *kafkaProducer) Close() {
	p.producer.Close()
}
