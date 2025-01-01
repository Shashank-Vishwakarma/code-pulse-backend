package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

func RateLimiter(noOfRequests int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		IPAddress, err := readUserIPAddress(c)
		if err != nil {
			logrus.Errorf("Error reading user IP address: RateLimiter Middleware: %v", err)
			response.HandleResponse(c, http.StatusNotFound, err.Error(), nil)
			c.Abort()
			return
		}

		cache := services.GetCache(IPAddress)
		if cache == nil { // first request
			err = services.SetCache(IPAddress, 1, duration)
			if err != nil {
				logrus.Errorf("Error setting cache: RateLimiter Middleware: %v", err)
				response.HandleResponse(c, http.StatusInternalServerError, err.Error(), nil)
				c.Abort()
				return
			}
		} else {
			cacheValue, err := strconv.Atoi(cache.(string))
			if err != nil {
				logrus.Errorf("Error converting cache value to int: RateLimiter Middleware: %v", err)
				response.HandleResponse(c, http.StatusInternalServerError, err.Error(), nil)
				c.Abort()
				return
			}

			if cacheValue >= noOfRequests {
				logrus.Error("Rate limit exceeded: RateLimiter Middleware")
				response.HandleResponse(c, http.StatusTooManyRequests, "Rate limit exceeded: Please try again after 1 hour", nil)
				c.Abort()
				return
			}

			err = services.SetCache(IPAddress, fmt.Sprintf("%d", cacheValue+1), duration)
			if err != nil {
				logrus.Errorf("Error setting cache: RateLimiter Middleware: %v", err)
				response.HandleResponse(c, http.StatusInternalServerError, err.Error(), nil)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
