package middlewares

import "github.com/gin-gonic/gin"

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
