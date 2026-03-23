package middleware

import "github.com/gin-gonic/gin"

func APIKey(required string) gin.HandlerFunc {
	if required == "" {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		provided := c.GetHeader("X-API-Key")
		if provided != required {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
