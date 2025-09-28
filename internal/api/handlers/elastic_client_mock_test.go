package handlers

import (
	"context"

	"github.com/rodrigogmartins/log-processor/internal/service"
)

type MockElasticSearchClient struct {
	IndexedLogs map[string]service.Log
	SearchFunc  func(ctx context.Context, index string, query map[string]interface{}, size int) ([]service.Log, error)
}

func NewMockElasticSearchClient() *MockElasticSearchClient {
	return &MockElasticSearchClient{
		IndexedLogs: make(map[string]service.Log),
	}
}

func (m *MockElasticSearchClient) Index(ctx context.Context, index string, id string, body interface{}) error {
	logEntry, ok := body.(service.Log)
	if !ok {
		return nil
	}
	m.IndexedLogs[id] = logEntry
	return nil
}

func (m *MockElasticSearchClient) SearchLogs(ctx context.Context, index string, query map[string]interface{}, size int) ([]service.Log, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, index, query, size)
	}

	logs := make([]service.Log, 0, len(m.IndexedLogs))
	for _, l := range m.IndexedLogs {
		logs = append(logs, l)
	}
	if len(logs) > size {
		logs = logs[:size]
	}
	return logs, nil
}
