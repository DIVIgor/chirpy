package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DIVIgor/chirpy/internal/database"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Create chirp by message and user id (for now)
func (cfg *apiConfig) handlerCreateChirp(writer http.ResponseWriter, req *http.Request) {
	type chirpPost struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	data := chirpPost{}
	err := decoder.Decode(&data)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't decode the provided parameters:", err)
		return
	}

	validatedBody, err := validateChirp(data.Body)
	if err != nil {
		respWithErr(writer, http.StatusBadRequest, err.Error(), err)
	}

	// save to DB
	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   validatedBody,
		UserID: data.UserID,
	})
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't save chirp to DB:", err)
		return
	}

	respJSON(writer, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}

// Retrieve a full list of chirps (for now)
func (cfg *apiConfig) handlerGetChirpList(writer http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		log.Println("Couldn't retrieve chirps:", err)
	}

	chirpList := []Chirp{}
	for _, chirp := range chirps {
		chirpList = append(chirpList, Chirp{
			ID:        chirp.ID,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
	}

	respJSON(writer, http.StatusOK, chirpList)
}

// Get a single chirp by its ID parsed from URL
func (cfg *apiConfig) handlerGetChirp(writer http.ResponseWriter, req *http.Request) {
	postID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respWithErr(writer, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), postID)
	if err != nil {
		respWithErr(writer, http.StatusNotFound, "Coudn't get chirp", err)
		return
	}

	respJSON(writer, http.StatusOK, Chirp{
		ID:        chirp.ID,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	})
}
