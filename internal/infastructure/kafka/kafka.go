package kafka

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var errUnknownType = errors.New("an unknown error was returned while receiving event from the channel")

const (
	flushtimeout = 5000
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(msg, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value:     []byte(msg),
		Timestamp: time.Now(),
	}

	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("failed to send message in channel %w", err)
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return fmt.Errorf("an kafka.error was returned while receiving event from the channel: %w", ev)
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushtimeout)
	p.producer.Close()
}
