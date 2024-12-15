package middlewares

import (
	"net/http"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(config.Config.JWT_TOKEN_COOKIE)
		if err != nil {
			logrus.Errorf("Error getting token: Authorization Middleware: %v", err)
			response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		if token == "" {
			logrus.Error("Unauthorized: Token not found: Authorization Middleware")
			response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Config.JWT_SECRET_KEY), nil
		})
		if err != nil {
			logrus.Errorf("Unauthorized: Invalid token: Authorization Middleware: %v", err)
			response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized: Invalid token", nil)
			c.Abort()
			return
		}

		if !jwtToken.Valid {
			logrus.Errorf("Unauthorized: Invalid token: Authorization Middleware")
			response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized: Invalid token", nil)
			c.Abort()
			return
		}

		// set the user data in the context
		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			logrus.Error("Error decoding the JWT claims: Authorization Middleware")
			response.HandleResponse(c, http.StatusUnauthorized, "Unauthorized: Invalid token", nil)
			c.Abort()
			return
		}

		userData := utils.JWTPayload{
			UserID:   claims["user_id"].(string),
			Name:     claims["name"].(string),
			Email:    claims["email"].(string),
			Username: claims["username"].(string),
		}
		c.Set(config.Config.JWT_DECODED_PAYLOAD, userData)

		c.Next()
	}
}
