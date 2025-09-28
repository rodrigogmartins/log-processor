package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/rodrigogmartins/log-processor/internal/service"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

type MockLogService struct {
	Processed  []service.Log
	ShouldFail bool
}

func (m *MockLogService) Process(ctx context.Context, logEntry service.Log) error {
	time.Sleep(10 * time.Millisecond)
	m.Processed = append(m.Processed, logEntry)
	return nil
}
func (m *MockLogService) SearchLogs(ctx context.Context, query map[string]interface{}, size int) ([]service.Log, error) {
	if m.Processed == nil {
		return []service.Log{}, nil
	}

	if size > len(m.Processed) {
		size = len(m.Processed)
	}
	return m.Processed[:size], nil
}

func TestProcessor_Start(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	messages := []kafka.Message{
		{Key: []byte("1"), Value: []byte("msg1")},
		{Key: []byte("2"), Value: []byte("msg2")},
		{Key: []byte("3"), Value: []byte("msg3")},
	}

	mockReader := &MockKafkaReader{Messages: messages}
	mockService := &MockLogService{}

	processor := NewProcessor(mockReader, mockService, 1, 1, time.Millisecond)

	done := make(chan struct{})
	go func() {
		_ = processor.Start(ctx)
		close(done)
	}()

	for {
		if len(mockService.Processed) == len(messages) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	cancel()
	<-done

	assert.Len(t, mockService.Processed, 3)
	assert.Len(t, mockReader.CommittedMsgs, 3)
	assert.True(t, mockReader.CloseCalled)
}

func TestProcessor_WorkerPool(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	messages := []kafka.Message{
		{Key: []byte("1"), Value: []byte("msg1")},
		{Key: []byte("2"), Value: []byte("msg2")},
		{Key: []byte("3"), Value: []byte("msg3")},
		{Key: []byte("4"), Value: []byte("msg4")},
		{Key: []byte("5"), Value: []byte("msg5")},
	}

	mockReader := &MockKafkaReader{Messages: messages}
	mockService := &MockLogService{}

	processor := NewProcessor(mockReader, mockService, 3, 2, 10*time.Millisecond)

	done := make(chan struct{})
	go func() {
		_ = processor.Start(ctx)
		close(done)
	}()

	timeout := time.After(5 * time.Second)
	for {
		if len(mockService.Processed) == len(messages) {
			break
		}
		select {
		case <-timeout:
			t.Fatal("timeout: not all messages processed")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	cancel()
	<-done

	assert.Len(t, mockService.Processed, 5)
	assert.Len(t, mockReader.CommittedMsgs, 5)
	assert.True(t, mockReader.CloseCalled)
}
