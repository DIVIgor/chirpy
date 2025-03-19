package main

import (
	"strings"
)

// Replace certain words with "****"
func censorWords(text string) (modText string) {
	replacement := "****"
	wordList := map[string]struct{}{  // perhaps the wordlist should be read from a textfile
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	words := strings.Split(text, " ")
	for idx, word := range words {
		if _, exists := wordList[strings.ToLower(word)]; exists {
			words[idx] = replacement
		}
	}

	return strings.Join(words, " ")
}