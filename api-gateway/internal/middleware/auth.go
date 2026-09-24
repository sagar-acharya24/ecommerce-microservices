package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const UserIDKey = "user_id"

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
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
				"error": "invalid authorization header",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(
			tokenString,
			&Claims{},
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(jwtSecret), nil
			},
		)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || claims.UserID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		// Store authenticated user ID in Gin context.
		c.Set(UserIDKey, claims.UserID)

		// Never trust a client-supplied X-User-ID header.
		c.Request.Header.Del("X-User-ID")

		// Add the user identity derived from the validated JWT.
		c.Request.Header.Set(
			"X-User-ID",
			strconv.FormatUint(uint64(claims.UserID), 10),
		)

		c.Next()
	}
}
