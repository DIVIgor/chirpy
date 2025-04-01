package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const tokenTypeAccess TokenType = "chirpy-access"
const Bearer string = "Bearer"
const ApiBearer string = "ApiKey"

// Hash password with Bcrypt
func HashPassword(password string) (hashedPW string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Couldn't hash password:", err)
		return
	}

	return string(hash), err
}

// Check hash
func CheckPasswordHash(password, hash string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("passwords don't match")
	}

	return err
}

// Create and sign JWT
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (signedToken string, err error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(tokenTypeAccess),
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
	})

	return token.SignedString(signingKey)
}

// Validate JWT by checking token, user, and issuer
func ValidateJWT(tokenString, tokenSecret string) (userID uuid.UUID, err error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return
	}
	if issuer != string(tokenTypeAccess) {
		return userID, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return userID, fmt.Errorf("invalid user ID: %w", err)
	}

	return id, err
}

// Check request headers for token and validate it. Return cleaned token string.
func GetBearerToken(headers http.Header, bearer string) (tokenStr string, err error) {
	token := headers.Get("Authorization")
	if len(token) == 0 {
		return tokenStr, errors.New("no auth header included in request")
	}

	splittedToken := strings.Split(token, " ")
	if len(splittedToken) < 2 || splittedToken[0] != bearer || len(splittedToken[1]) == 0 {
		return tokenStr, errors.New("wrong token format")
	}

	return splittedToken[1], err
}

// Create random 256-bit refresh token encoded in hex
func MakeRefreshToken() (refreshToken string, err error) {
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		return refreshToken, err
	}

	return hex.EncodeToString(key), err
}
