package auth

import (
	"testing"
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
