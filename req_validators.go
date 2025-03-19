package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func handlerChirpValidation(writer http.ResponseWriter, req *http.Request) {
	type request struct {
		Body string `json:"body"`
	}
	type cleanResp struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := request{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Couldn't decode the provided parameters: %s", err)
		respWithErr(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	// limit messages to 140 symbols
	if len(params.Body) > 140 {
		respWithErr(writer, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// censor certain words
	// since length of the smallest word to censor is 5 chars
	if len(params.Body) > 5 {
		params.Body = censorWords(params.Body)
	}

	respJSON(writer, http.StatusOK, cleanResp{CleanedBody: params.Body})
}