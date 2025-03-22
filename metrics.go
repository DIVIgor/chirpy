package main

import (
	"fmt"
	"net/http"
)

// Check the number of requests to server
func (cfg *apiConfig) handlerCountVisits(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
}

// Reset the number of requests to server and users
func (cfg *apiConfig) handlerResetVisits(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// check if the app in dev stage
	if cfg.platform != "dev" {
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("Forbidden"))
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Visitor counter has been reset."))

	cfg.fileserverHits.Store(0)
	cfg.dbQueries.ClearUsers(req.Context())
}
