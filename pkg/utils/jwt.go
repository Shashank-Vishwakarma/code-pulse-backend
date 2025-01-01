package utils

import (
	"net/http"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTPayload struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func GenerateToken(payload JWTPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       payload.ID,
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

func GetDecodedUserFromContext(c *gin.Context) (JWTPayload, error) {
	userData, exists := c.Get(config.Config.JWT_DECODED_PAYLOAD)
	if !exists {
		response.HandleResponse(c, http.StatusBadRequest, "User data not found", nil)
		return JWTPayload{}, nil
	}

	decodedUser, ok := userData.(JWTPayload)
	if !ok {
		response.HandleResponse(c, http.StatusBadRequest, "Invalid user data", nil)
		return JWTPayload{}, nil
	}

	return decodedUser, nil
}
