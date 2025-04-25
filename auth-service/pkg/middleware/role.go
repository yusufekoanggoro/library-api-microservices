package middleware

import (
	"auth-service/pkg/shared/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *middleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Role not found")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "Invalid Role type")
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "You do not have permission")
		c.Abort()
	}
}
