package service

import (
	"context"
)

type MockElasticSearchClient struct {
	IndexedLogs map[string]Log
	SearchFunc  func(ctx context.Context, index string, query map[string]interface{}, size int) ([]Log, error)
}

func NewMockElasticSearchClient() *MockElasticSearchClient {
	return &MockElasticSearchClient{
		IndexedLogs: make(map[string]Log),
	}
}

func (m *MockElasticSearchClient) Index(ctx context.Context, index string, id string, body interface{}) error {
	logEntry, ok := body.(Log)
	if !ok {
		return nil
	}
	m.IndexedLogs[id] = logEntry
	return nil
}

func (m *MockElasticSearchClient) SearchLogs(ctx context.Context, index string, query map[string]interface{}, size int) ([]Log, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, index, query, size)
	}

	logs := make([]Log, 0, len(m.IndexedLogs))
	for _, l := range m.IndexedLogs {
		logs = append(logs, l)
	}
	if len(logs) > size {
		logs = logs[:size]
	}
	return logs, nil
}
