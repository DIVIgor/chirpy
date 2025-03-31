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

// Parsable user model for CRUD operations
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Success response structure
type response struct {
	User
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// Read and decode request body
func decodeRequest(req *http.Request) (data userAuth, err error) {
	decoder := json.NewDecoder(req.Body)
	data = userAuth{}
	err = decoder.Decode(&data)
	if err != nil {
		return data, fmt.Errorf("couldn't decode parameters: %w", err)
	}

	return data, err
}

// Get and check refresh token from headers and generate an access token
func (cfg *apiConfig) handlerRefreshAccess(writer http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	// get refresh token from headers and check refresh token format
	reqToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respWithErr(writer, http.StatusBadRequest, "Couldn't find token", err)
		return
	}
	// check user
	user, err := cfg.dbQueries.GetUserFromToken(req.Context(), reqToken)
	if err != nil {
		respWithErr(writer, http.StatusUnauthorized, "Couldn't find user for refresh", err)
		return
	}
	// generate new access token
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't create JWT", err)
		return
	}

	respJSON(writer, http.StatusOK, response{
		Token: accessToken,
	})
}

// Get and check refresh token from headers and mark it as revoked in DB
func (cfg *apiConfig) handlerRevokeAccess(writer http.ResponseWriter, req *http.Request) {
	reqToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respWithErr(writer, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}
	err = cfg.dbQueries.RevokeToken(req.Context(), reqToken)
	if err != nil {
		respWithErr(writer, http.StatusUnauthorized, "Couldn't revoke session", err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

// Create user with email and password
func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	data, err := decodeRequest(req)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, err.Error(), err)
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          data.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't create user", err)
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

// Update email and/or password with provided credentials and valid token
func (cfg *apiConfig) handlerUpdateUser(writer http.ResponseWriter, req *http.Request) {
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

	data, err := decodeRequest(req)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, err.Error(), err)
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          data.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't update credentials", err)
		return
	}

	respJSON(writer, http.StatusOK, response{
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
	data, err := decodeRequest(req)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, err.Error(), err)
		return
	}

	user, err := cfg.dbQueries.GetUser(req.Context(), data.Email)
	if err != nil {
		respWithErr(writer, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	err = auth.CheckPasswordHash(data.Password, user.HashedPassword)
	if err != nil {
		respWithErr(writer, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	// Refresh Token should expire in 60 days
	savedToken, err := cfg.dbQueries.SaveRefreshToken(req.Context(), database.SaveRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respJSON(writer, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token:        accessToken,
		RefreshToken: savedToken.Token,
	})
}
