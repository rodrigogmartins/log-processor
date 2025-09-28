package kafka

import (
	"context"
	"io"
	"time"

	"github.com/segmentio/kafka-go"
)

type MockKafkaReader struct {
	Messages      []kafka.Message
	Index         int
	CommittedMsgs []kafka.Message
	CloseCalled   bool
}

func (m *MockKafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	if m.Index >= len(m.Messages) {
		// simula loop do Kafka sem travar o teste
		select {
		case <-ctx.Done():
			return kafka.Message{}, ctx.Err()
		case <-time.After(5 * time.Millisecond):
			return kafka.Message{}, io.EOF
		}
	}

	msg := m.Messages[m.Index]
	m.Index++
	return msg, nil
}

func (m *MockKafkaReader) CommitMessage(msg kafka.Message) error {
	m.CommittedMsgs = append(m.CommittedMsgs, msg)
	return nil
}

func (m *MockKafkaReader) Close() error {
	m.CloseCalled = true
	return nil
}
