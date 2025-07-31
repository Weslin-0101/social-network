package security

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func IsStrongPassword(password string) bool {
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false
	specialChars := "!@#$%^&*()-_=+[]{}|;:',.<>?/~`"

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		}
		if char >= 'a' && char <= 'z' {
			hasLower = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if strings.ContainsRune(specialChars, char) {
			hasSpecial = true
		}

		if hasUpper && hasLower && hasNumber && hasSpecial {
			return true
		}
	}

	return false
}