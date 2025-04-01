package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/DIVIgor/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Check server status
func handlerReadiness(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	// .env params
	platform  string // dev or prod
	jwtSecret string
	polkaKey  string
}

// Count requests to the server (main paths only)
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(writer, req)
	})
}

// DRY paths by using in helper functions with prefilled `section`
//
// Keeps HTTP method in upper case and adds space after it, if provided.
//
// `path` parameter should start with "/"
func pathMod(reqMethod, section, path string) (modifiedPath string) {
	if len(reqMethod) > 2 {
		reqMethod = strings.ToUpper(reqMethod)
		reqMethod = strings.TrimSpace(reqMethod) + " "
	}
	return fmt.Sprintf("%s%s%s", reqMethod, section, path)
}

// API path modifier
func apiPath(reqMethod, path string) (modifiedPath string) {
	return pathMod(reqMethod, "/api", path)
}

// Admin path modifier
func adminPath(reqMethod, path string) (modifiedPath string) {
	return pathMod(reqMethod, "/admin", path)
}

func main() {
	const filePathRoot string = "."
	const port string = "8080"

	// get DB path and load it
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	jwtSecret := os.Getenv("SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT secret is not set.")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("Polka key is not set.")
	}

	apiCfg := &apiConfig{
		dbQueries: database.New(db),
		platform:  os.Getenv("PLATFORM"),
		jwtSecret: jwtSecret,
		polkaKey:  polkaKey,
	}
	mux := http.NewServeMux()

	// Main path
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))
	// Secondary paths:
	// • API:
	// 	- server health
	mux.HandleFunc(apiPath("GET", "/healthz"), handlerReadiness)
	// 	- account
	mux.HandleFunc(apiPath("POST", "/users"), apiCfg.handlerCreateUser)
	mux.HandleFunc(apiPath("PUT", "/users"), apiCfg.handlerUpdateUser)
	mux.HandleFunc(apiPath("POST", "/login"), apiCfg.handlerLogin)
	mux.HandleFunc(apiPath("POST", "/refresh"), apiCfg.handlerRefreshAccess)
	mux.HandleFunc(apiPath("POST", "/revoke"), apiCfg.handlerRevokeAccess)
	// 	- posts
	mux.HandleFunc(apiPath("POST", "/chirps"), apiCfg.handlerCreateChirp)
	mux.HandleFunc(apiPath("GET", "/chirps"), apiCfg.handlerGetChirpList)
	mux.HandleFunc(apiPath("GET", "/chirps/{chirpID}"), apiCfg.handlerGetChirp)
	mux.HandleFunc(apiPath("DELETE", "/chirps/{chirpID}"), apiCfg.handlerDeleteChirp)
	// 	- webhooks
	mux.HandleFunc(apiPath("POST", "/polka/webhooks"), apiCfg.handlerUpgradeUserPlan)
	// • Administration:
	// 	- metrics
	mux.HandleFunc(adminPath("GET", "/metrics"), apiCfg.handlerCountVisits)
	// 	- reset DB
	mux.HandleFunc(adminPath("POST", "/reset"), apiCfg.handlerResetVisits)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Simple info
	log.Println("Serving files from", filePathRoot, "on port:", port)

	log.Fatal(server.ListenAndServe())
}
