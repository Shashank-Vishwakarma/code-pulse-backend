package utils

import (
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTPayload struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func GenerateToken(payload JWTPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  payload.UserID,
		"name":     payload.Name,
		"email":    payload.Email,
		"username": payload.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Config.JWT_SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
