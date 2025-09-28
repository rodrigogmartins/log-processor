package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rodrigogmartins/log-processor/internal/service"
)

type LogHandler struct {
	LogService service.ElasticSearchClient
	Index      string
}

// GET /logs
func (h *LogHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	query := map[string]interface{}{
		"match_all": map[string]interface{}{},
	}
	size := 100

	logs, err := h.LogService.SearchLogs(r.Context(), h.Index, query, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(logs)
}

// GET /logs?level=INFO
func (h *LogHandler) ListLogsByLevel(w http.ResponseWriter, r *http.Request) {
	level := r.URL.Query().Get("level")
	if level == "" {
		h.ListLogs(w, r)
		return
	}

	query := map[string]interface{}{
		"match": map[string]interface{}{
			"level": level,
		},
	}
	size := 100

	logs, err := h.LogService.SearchLogs(r.Context(), h.Index, query, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(logs)
}

// GET /logs/{id}
func (h *LogHandler) GetLogByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	query := map[string]interface{}{
		"term": map[string]interface{}{
			"id": id,
		},
	}
	size := 1

	logs, err := h.LogService.SearchLogs(r.Context(), h.Index, query, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(logs) == 0 {
		http.Error(w, "log not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(logs[0])
}
