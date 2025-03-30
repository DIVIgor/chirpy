package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// Password examples
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name        string
		password    string
		hash        string
		expectedErr bool
	}{
		{
			name:        "Correct password",
			password:    password1,
			hash:        hash1,
			expectedErr: false,
		},
		{
			name:        "Incorrect password",
			password:    "wrongPassword",
			hash:        hash1,
			expectedErr: true,
		},
		{
			name:        "Password doesn't match different hash",
			password:    password1,
			hash:        hash2,
			expectedErr: true,
		},
		{
			name:        "Empty password",
			password:    "",
			hash:        hash1,
			expectedErr: true,
		},
		{
			name:        "Invalid hash",
			password:    password1,
			hash:        "invalidhash",
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := CheckPasswordHash(testCase.password, testCase.hash)
			if (err != nil) != testCase.expectedErr {
				t.Errorf("CheckPasswordHash() error = %v, expectedErr %v", err, testCase.expectedErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(testCase.tokenString, testCase.tokenSecret)
			if (err != nil) != testCase.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if gotUserID != testCase.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, testCase.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	headerName := "Authorization"
	token, _ := MakeJWT(uuid.New(), "secret", time.Hour)
	fullToken := fmt.Sprintf("%s %s", Bearer, token)
	validHeaders := http.Header{headerName: []string{fullToken}}

	tests := []struct {
		name        string
		headers     http.Header
		token       string
		expectedErr bool
	}{
		{
			name:        "Valid bearer token",
			headers:     validHeaders,
			token:       token,
			expectedErr: false,
		},
		{
			name:        "No authorization header",
			headers:     http.Header{"Auth": []string{fullToken}},
			token:       "",
			expectedErr: true,
		},
		{
			name:        "Empty string",
			headers:     http.Header{headerName: []string{""}},
			token:       "",
			expectedErr: true,
		},
		{
			name:        "No bearer prefix",
			headers:     http.Header{headerName: []string{" " + token}},
			token:       "",
			expectedErr: true,
		},
		{
			name:        "Token string without spaces",
			headers:     http.Header{headerName: []string{fmt.Sprintf("%s%s", Bearer, token)}},
			token:       "",
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tokenString, err := GetBearerToken(testCase.headers)
			if (err != nil) != testCase.expectedErr {
				t.Errorf("GetBearerToken() error = %v, expectedErr %v", err, testCase.expectedErr)
				return
			}
			if tokenString != testCase.token {
				t.Errorf("GetBearerToken() token = %v, expectedToken %v", tokenString, token)
			}
		})
	}
}
