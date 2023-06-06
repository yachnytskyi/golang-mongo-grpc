package utility

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Encrypt the password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash the password %w", err)
	}

	return string(hashedPassword), nil
}
