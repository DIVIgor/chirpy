package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// simple user creation (requires proper auth)
func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	// valid request body
	type request struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	data := request{}
	err := decoder.Decode(&data)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't decode parameters:", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), data.Email)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't create user:", err)
		return
	}

	respJSON(writer, http.StatusCreated, User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
