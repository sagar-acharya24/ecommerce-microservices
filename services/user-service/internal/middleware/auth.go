package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sagar-acharya24/ecommerce-microservices/services/user-service/internal/auth"
)

const UserIDKey = "user_id"

func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header must use Bearer token",
			})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)

		c.Next()
	}
}
