package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader *kafka.Reader
	Topic  string
}

func NewKafkaConsumer(brokers []string, topic string, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		Topic: topic,
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			GroupID:     groupID,
			GroupTopics: []string{topic},
			StartOffset: kafka.FirstOffset,
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
			MaxWait:     500 * time.Millisecond,
		}),
	}
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.Reader.ReadMessage(ctx)
}

func (c *KafkaConsumer) CommitMessage(msg kafka.Message) error {
	return c.Reader.CommitMessages(context.Background(), msg)
}

func (c *KafkaConsumer) Close() error {
	return c.Reader.Close()
}
