package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DIVIgor/chirpy/internal/auth"
	"github.com/DIVIgor/chirpy/internal/database"
	"github.com/google/uuid"
)

// Valid request body with credentials
type userAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Read and decode request body
func decodeCredentials(req *http.Request) (data userAuth, err error) {
	decoder := json.NewDecoder(req.Body)
	data = userAuth{}
	err = decoder.Decode(&data)
	if err != nil {
		return data, fmt.Errorf("couldn't decode parameters: %w", err)
	}

	return data, err
}

// Success response structure
type response struct {
	User
}

// Parsable user model for CRUD operations
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Create user with email and password
func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	data, err := decodeCredentials(req)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, err.Error(), err)
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't hash password:", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          data.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't create user:", err)
		return
	}

	respJSON(writer, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// Login with email and password
func (cfg *apiConfig) handlerLogin(writer http.ResponseWriter, req *http.Request) {
	data, err := decodeCredentials(req)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, err.Error(), err)
		return
	}

	user, usrErr := cfg.dbQueries.GetUser(req.Context(), data.Email)
	pwdErr := auth.CheckPasswordHash(data.Password, user.HashedPassword)
	if usrErr != nil || pwdErr != nil {
		respWithErr(writer, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	respJSON(writer, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
