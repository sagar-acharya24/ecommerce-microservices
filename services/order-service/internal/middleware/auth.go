package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func AuthenticatedUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDHeader := c.GetHeader("X-User-ID")

		if userIDHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authenticated user identity is required",
			})
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(userIDHeader, 10, 64)
		if err != nil || userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authenticated user identity",
			})
			c.Abort()
			return
		}

		c.Set(UserIDKey, uint(userID))

		c.Next()
	}
}
