package authentication

import (
	"backend/src/config"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/openpgp/errors"
)

func GenerateToken(userID uint64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 24).Unix()
	permissions["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return token.SignedString([]byte(config.SecretKey))
}

func ValidateToken(r *http.Request) error {
	tokenString := extractToken(r)
	
	token, err := jwt.Parse(tokenString, getVerificationKey)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.InvalidArgumentError("invalid token")
}

func ExtractUserID(r *http.Request) (uint64, error) {
	tokenString := extractToken(r)

	token, err := jwt.Parse(tokenString, getVerificationKey)
	if err != nil {
		return 0, err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIDFloat, ok := permissions["user_id"].(float64); ok {
			return uint64(userIDFloat), nil
		}

		if userIDInt, ok := permissions["user_id"].(int64); ok {
			return uint64(userIDInt), nil
		}

		return 0, errors.InvalidArgumentError("invalid user_id type in token")
	}

	return 0, errors.InvalidArgumentError("invalid token claims")
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func getVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}