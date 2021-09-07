package eventbus

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/ozonva/ova-checklist-api/internal/config"
)

// kafkaEventBus implements EventBus
type kafkaEventBus struct {
	writer *kafka.Writer
}

func NewEventBusOverKafka(cfg *config.KafkaConfig) EventBus {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	writer := &kafka.Writer{
		Addr:     kafka.TCP(address),
		Topic:    cfg.Topic,
		Balancer: kafka.Murmur2Balancer{},
	}
	return &kafkaEventBus{
		writer: writer,
	}
}

func (k *kafkaEventBus) Send(ctx context.Context, events ...Event) error {
	messages := make([]kafka.Message, 0, len(events))
	for _, event := range events {
		messages = append(messages, kafka.Message{
			Key:   []byte(event.Key),
			Value: event.Value,
		})
	}
	return k.writer.WriteMessages(ctx, messages...)
}

func (k *kafkaEventBus) Close() error {
	if k.writer != nil {
		if err := k.writer.Close(); err != nil {
			return err
		}
	}
	return nil
}
