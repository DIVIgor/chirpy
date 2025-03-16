package main

import (
	"log"
	"net/http"
)


func main() {
	const filePathRoot string = "."
	const port string = "8080"
	
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))

	mux.HandleFunc("/healthz", handlerReadiness)

	server := &http.Server {
		Addr: ":" + port,
		Handler: mux,
	}

	log.Println("Serving files from", filePathRoot, "on port:", port)
	
	log.Fatal(server.ListenAndServe())
}

// Check server status
func handlerReadiness(writer http.ResponseWriter, req *http.Request){
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}