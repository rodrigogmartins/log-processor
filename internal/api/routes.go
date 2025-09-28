package api

import (
	"github.com/gorilla/mux"
	"github.com/rodrigogmartins/log-processor/internal/api/handlers"
	"github.com/rodrigogmartins/log-processor/internal/service"
)

func NewRouter(logService service.ElasticSearchClient, index string) *mux.Router {
	handler := &handlers.LogHandler{
		LogService: logService,
		Index:      index,
	}

	r := mux.NewRouter()
	r.HandleFunc("/logs", handler.ListLogs).Methods("GET")
	r.HandleFunc("/logs/by-level", handler.ListLogsByLevel).Methods("GET")
	r.HandleFunc("/logs/{id}", handler.GetLogByID).Methods("GET")

	return r
}
