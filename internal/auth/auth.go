package auth

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

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
