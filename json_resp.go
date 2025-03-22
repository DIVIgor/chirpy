package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Form, log, and send request/response error
func respWithErr(writer http.ResponseWriter, statusCode int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if statusCode > 499 {
		log.Printf("Responding with %d server error: %s", statusCode, msg)
	}

	type errResp struct {
		Err string `json:"error"`
	}
	respJSON(writer, statusCode, errResp{Err: msg})
}

// Form and send JSON response
func respJSON(writer http.ResponseWriter, statusCode int, payload interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(statusCode)
	writer.Write(resp)
}
