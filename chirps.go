package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/DIVIgor/chirpy/internal/auth"
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
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respWithErr(writer, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respWithErr(writer, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	type chirpPost struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	data := chirpPost{}
	err = decoder.Decode(&data)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't decode parameters:", err)
		return
	}

	validatedBody, err := validateChirp(data.Body)
	if err != nil {
		respWithErr(writer, http.StatusBadRequest, err.Error(), err)
	}

	// save to DB
	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   validatedBody,
		UserID: userID,
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

// Delete chirp by its ID and owner ID
func (cfg *apiConfig) handlerDeleteChirp(writer http.ResponseWriter, req *http.Request) {
	postID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respWithErr(writer, http.StatusNotFound, "Chirp not found", err)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respWithErr(writer, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respWithErr(writer, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	post, err := cfg.dbQueries.DeleteChirp(req.Context(), database.DeleteChirpParams{
		ID:     postID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respWithErr(writer, http.StatusNotFound, "Chirp not found", err)
		} else {
			respWithErr(writer, http.StatusInternalServerError, "Couldn't delete chirp", err)
		}
		return
	}
	// empty post means that user is not the owner of this chirp
	if post == (database.Chirp{}) {
		respWithErr(writer, http.StatusForbidden, "You can't delete this chirp", err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
