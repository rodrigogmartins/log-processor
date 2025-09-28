package db

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/rodrigogmartins/log-processor/internal/service"

	esv8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
)

type MockTransport struct {
	Response *http.Response
	Err      error
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestElasticSearchClient_Index(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: 201,
		Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		Body:       io.NopCloser(bytes.NewBufferString(`{"result":"created"}`)),
	}

	cfg := esv8.Config{
		Transport: &MockTransport{Response: mockResp},
	}
	client, _ := esv8.NewClient(cfg)
	esClient := &ElasticSearchClient{Client: client}

	logEntry := service.Log{
		ID:      "1",
		Message: "Test log",
		Level:   "INFO",
	}

	err := esClient.Index(context.Background(), "logs-index", logEntry.ID, logEntry)
	assert.NoError(t, err)
}

func TestElasticSearchClient_SearchLogs(t *testing.T) {
	respJSON := `{
		"hits": {
			"hits": [
				{"_source": {"id":"1","message":"msg1","level":"INFO"}},
				{"_source": {"id":"2","message":"msg2","level":"INFO"}}
			]
		}
	}`

	mockResp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		Body:       io.NopCloser(bytes.NewBufferString(respJSON)),
	}

	cfg := esv8.Config{
		Transport: &MockTransport{Response: mockResp},
	}
	client, _ := esv8.NewClient(cfg)
	esClient := &ElasticSearchClient{Client: client}

	query := map[string]interface{}{
		"match": map[string]interface{}{"level": "INFO"},
	}

	logs, err := esClient.SearchLogs(context.Background(), "logs-index", query, 10)
	assert.NoError(t, err)
	assert.Len(t, logs, 2)
	assert.Equal(t, "msg1", logs[0].Message)
	assert.Equal(t, "msg2", logs[1].Message)
}
