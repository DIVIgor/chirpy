package main

import "errors"

func validateChirp(message string) (string, error) {
	// limit messages to 140 symbols
	if len(message) > 140 {
		return "", errors.New("Chirp is too long")
	}

	// censor certain words
	// since length of the smallest word to censor is 5 chars
	if len(message) > 5 {
		message = censorWords(message)
	}

	return message, nil
}
