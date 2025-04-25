package middleware

import (
	"auth-service/pkg/token"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *middleware) JWTAuthMiddleware(token token.Token) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenClaims, err := token.ValidateToken(parts[1])
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userId", tokenClaims.UserID)
		c.Set("role", tokenClaims.Role)
		c.Next()
	}
}
