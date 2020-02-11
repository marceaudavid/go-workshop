package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash ...
func Hash(password string, key string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password+key), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Compare ...
func Compare(hash string, password string, key string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+key))
	if err != nil {
		return false
	}
	return true
}
