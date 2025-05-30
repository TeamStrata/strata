package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// Return a bcrypt password hash, return empty string if error occurred.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Check if login password matches user's hashed password.
func AuthenticateUser(userPassword string, loginPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(loginPassword))
}
