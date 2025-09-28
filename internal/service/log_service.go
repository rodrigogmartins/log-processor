package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type ElasticSearchClient interface {
	Index(ctx context.Context, index string, id string, body interface{}) error
	SearchLogs(ctx context.Context, index string, query map[string]interface{}, size int) ([]Log, error)
}

type LogServiceInterface interface {
	Process(ctx context.Context, logEntry Log) error
	SearchLogs(ctx context.Context, query map[string]interface{}, size int) ([]Log, error)
}

type LogService struct {
	esClient ElasticSearchClient
	index    string
}

type Log struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
}

func NewLogService(esClient ElasticSearchClient, index string) *LogService {
	return &LogService{
		esClient: esClient,
		index:    index,
	}
}

func (s *LogService) Process(ctx context.Context, logEntry Log) error {
	if logEntry.ID == "" || logEntry.Message == "" {
		return errors.New("invalid log: empty ID or message")
	}

	if logEntry.Timestamp.IsZero() {
		logEntry.Timestamp = time.Now().UTC()
	}

	if err := s.esClient.Index(ctx, s.index, logEntry.ID, logEntry); err != nil {
		return err
	}

	return nil
}

func (s *LogService) SearchLogs(ctx context.Context, query map[string]interface{}, size int) ([]Log, error) {
	return s.esClient.SearchLogs(ctx, s.index, query, size)
}

func EncodeLog(logEntry Log) ([]byte, error) {
	return json.Marshal(logEntry)
}
