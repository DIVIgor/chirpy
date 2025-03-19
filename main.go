package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

// Check server status
func handlerReadiness(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}


type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(writer, req)
	})
}

// path modifier for API calls
func apiPath(reqMethod, path string) (modifiedPath string) {
	if len(reqMethod) > 2 {
		reqMethod = strings.ToUpper(reqMethod)
		reqMethod = strings.TrimSpace(reqMethod) + " "
	}
	return fmt.Sprintf("%s/api%s", reqMethod, path)
}


func main() {
	const filePathRoot string = "."
	const port string = "8080"

	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))

	mux.HandleFunc(apiPath("GET", "/healthz"), handlerReadiness)
	mux.HandleFunc(apiPath("POST", "/validate_chirp"), handlerChirpValidation)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerCountVisits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetVisits)

	server := &http.Server {
		Addr: ":" + port,
		Handler: mux,
	}

	log.Println("Serving files from", filePathRoot, "on port:", port)

	log.Fatal(server.ListenAndServe())
}