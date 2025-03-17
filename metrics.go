package main

import (
	"fmt"
	"net/http"
)

// Check the number of requests to server
func (cfg *apiConfig) handlerCountVisits(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

// Reset the number of requests to server
func (cfg *apiConfig) handlerResetVisits(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Visitor counter has been reset."))

	cfg.fileserverHits.Store(0)
}