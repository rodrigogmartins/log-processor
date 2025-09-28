package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockElasticSearch struct {
	Indexed    []Log
	SearchFunc func(ctx context.Context, index string, query map[string]interface{}, size int) ([]Log, error)
	Err        error
}

func (m *MockElasticSearch) Index(ctx context.Context, index string, id string, body interface{}) error {
	if m.Err != nil {
		return m.Err
	}

	logEntry, ok := body.(Log)
	if !ok {
		return errors.New("invalid type")
	}

	m.Indexed = append(m.Indexed, logEntry)

	return nil
}

func (m *MockElasticSearch) SearchLogs(ctx context.Context, index string, query map[string]interface{}, size int) ([]Log, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, index, query, size)
	}

	logs := append([]Log{}, m.Indexed...)

	if len(logs) > size {
		logs = logs[:size]
	}
	return logs, nil
}

func TestProcess(t *testing.T) {

	t.Run("GIVEN valid log WHEN call Process THEN return nil", func(t *testing.T) {
		mockES := &MockElasticSearch{}
		service := NewLogService(mockES, "logs-index")
		logEntry := Log{
			ID:      "1",
			Level:   "INFO",
			Message: "Test log message",
			Source:  "unit-test",
		}

		err := service.Process(context.Background(), logEntry)

		assert.NoError(t, err, "must not return an execution error")
		assert.Len(t, mockES.Indexed, 1, "must have 1 indexed log")
		assert.Equal(t, "1", mockES.Indexed[0].ID)
		assert.Equal(t, "Test log message", mockES.Indexed[0].Message)
	})

	t.Run("GIVEN invalid log WHEN call Process THEN return error", func(t *testing.T) {
		mockES := &MockElasticSearch{}
		service := NewLogService(mockES, "logs-index")
		logEntry := Log{
			ID:      "",
			Level:   "INFO",
			Message: "",
		}

		err := service.Process(context.Background(), logEntry)

		assert.Error(t, err, "must return error for invalid log")
	})

	t.Run("GIVEN invalid log WHEN call Process THEN return error", func(t *testing.T) {
		mockES := &MockElasticSearch{
			Err: errors.New("ES down"),
		}
		service := NewLogService(mockES, "logs-index")
		logEntry := Log{ID: "1", Level: "INFO"}

		err := service.Process(context.Background(), logEntry)

		assert.Error(t, err, "must return ElasticSearch error")
	})
}

func TestLogService_SearchLogs(t *testing.T) {
	ctx := context.Background()
	mockES := NewMockElasticSearchClient()
	logService := NewLogService(mockES, "logs-index")

	// indexa logs
	for i := 1; i <= 5; i++ {
		logEntry := Log{
			ID:        string(rune(i)),
			Message:   "msg" + string(rune(i)),
			Level:     "INFO",
			Source:    "unit-test",
			Timestamp: time.Now().UTC(),
		}
		logService.Process(ctx, logEntry)
	}

	// realiza busca
	query := map[string]interface{}{
		"match": map[string]interface{}{"level": "INFO"},
	}
	logs, err := logService.SearchLogs(ctx, query, 10)
	assert.NoError(t, err)
	assert.Len(t, logs, 5)
}
