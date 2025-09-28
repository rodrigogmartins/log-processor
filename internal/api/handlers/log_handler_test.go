package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rodrigogmartins/log-processor/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestLogHandler_ListLogs(t *testing.T) {
	mockClient := &MockElasticSearchClient{
		SearchFunc: func(ctx context.Context, index string, query map[string]interface{}, size int) ([]service.Log, error) {
			return []service.Log{
				{ID: "1", Message: "log1", Timestamp: time.Now()},
				{ID: "2", Message: "log2", Timestamp: time.Now()},
			}, nil
		},
	}
	handler := &LogHandler{LogService: mockClient, Index: "logs-index"}

	req := httptest.NewRequest(http.MethodGet, "/logs", nil)
	w := httptest.NewRecorder()

	handler.ListLogs(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var logs []service.Log
	err := json.NewDecoder(resp.Body).Decode(&logs)
	assert.NoError(t, err)
	assert.Len(t, logs, 2)
}

func TestLogHandler_ListLogsByLevel(t *testing.T) {
	mockClient := &MockElasticSearchClient{
		SearchFunc: func(ctx context.Context, index string, query map[string]interface{}, size int) ([]service.Log, error) {
			if level, ok := query["match"].(map[string]interface{})["level"]; ok && level == "INFO" {
				return []service.Log{
					{ID: "1", Level: "INFO", Message: "log info", Timestamp: time.Now()},
				}, nil
			}
			return []service.Log{}, nil
		},
	}
	handler := &LogHandler{LogService: mockClient, Index: "logs-index"}

	req := httptest.NewRequest(http.MethodGet, "/logs?level=INFO", nil)
	w := httptest.NewRecorder()

	handler.ListLogsByLevel(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var logs []service.Log
	err := json.NewDecoder(resp.Body).Decode(&logs)
	assert.NoError(t, err)
	assert.Len(t, logs, 1)
	assert.Equal(t, "INFO", logs[0].Level)
}

func TestLogHandler_GetLogByID(t *testing.T) {
	mockClient := &MockElasticSearchClient{
		SearchFunc: func(ctx context.Context, index string, query map[string]interface{}, size int) ([]service.Log, error) {
			if term, ok := query["term"].(map[string]interface{})["id"]; ok && term == "1" {
				return []service.Log{
					{ID: "1", Message: "log1", Timestamp: time.Now()},
				}, nil
			}
			return []service.Log{}, nil
		},
	}
	handler := &LogHandler{LogService: mockClient, Index: "logs-index"}

	// Test ID existente
	req := httptest.NewRequest(http.MethodGet, "/logs?id=1", nil)
	w := httptest.NewRecorder()
	handler.GetLogByID(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var log service.Log
	err := json.NewDecoder(resp.Body).Decode(&log)
	assert.NoError(t, err)
	assert.Equal(t, "1", log.ID)

	// Test ID inexistente
	req2 := httptest.NewRequest(http.MethodGet, "/logs?id=999", nil)
	w2 := httptest.NewRecorder()
	handler.GetLogByID(w2, req2)
	resp2 := w2.Result()
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp2.StatusCode)

	// Test sem ID
	req3 := httptest.NewRequest(http.MethodGet, "/logs", nil)
	w3 := httptest.NewRecorder()
	handler.GetLogByID(w3, req3)
	resp3 := w3.Result()
	defer resp3.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp3.StatusCode)
}

func TestLogHandler_ErrorFromClient(t *testing.T) {
	mockClient := &MockElasticSearchClient{
		SearchFunc: func(ctx context.Context, index string, query map[string]interface{}, size int) ([]service.Log, error) {
			return nil, errors.New("client error")
		},
	}
	handler := &LogHandler{LogService: mockClient, Index: "logs-index"}

	req := httptest.NewRequest(http.MethodGet, "/logs", nil)
	w := httptest.NewRecorder()
	handler.ListLogs(w, req)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
