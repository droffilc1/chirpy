package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}

	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := badWordReplacement(params.Body, badWords)

	respondWithJSON(w, http.StatusOK, responseBody{
		CleanedBody: cleaned,
	})
}

func badWordReplacement(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, w := range words {
		if _, ok := badWords[strings.ToLower(w)]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
