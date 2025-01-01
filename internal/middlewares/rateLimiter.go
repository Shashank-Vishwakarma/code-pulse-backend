package middlewares

import (
	"fmt"
	"net/http"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/services"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func readUserIPAddress(c *gin.Context) (string, error) {
	IPAddress := c.Request.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = c.Request.Header.Get("X-Forwarded-For")
	}

	if IPAddress == "" {
		IPAddress = c.RemoteIP()
	}

	if IPAddress == "" {
		return "", fmt.Errorf("ip address not found")
	}

	return IPAddress, nil
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		IPAddress, err := readUserIPAddress(c)
		if err != nil {
			logrus.Errorf("Error reading user IP address: RateLimiter Middleware: %v", err)
			response.HandleResponse(c, http.StatusNotFound, err.Error(), nil)
			c.Abort()
			return
		}

		cache, err := services.GetCache(IPAddress)
		if err != nil {
			logrus.Errorf("Error getting IP Address from cache: RateLimiter Middleware: %v", err)
			response.HandleResponse(c, http.StatusInternalServerError, err.Error(), nil)
			c.Abort()
			return
		}

		cacheValue, ok := cache.(int)
		if !ok {
			logrus.Error("Could not convert cache value to int: RateLimiter Middleware")
			response.HandleResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", nil)
			c.Abort()
			return
		}

		if cacheValue >= 30 {
			logrus.Error("Rate limit exceeded: RateLimiter Middleware")
			response.HandleResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
