package middleware

import (
	"book-service/config"
	"book-service/pkg/token"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	BasicAuthMiddleware(cfg config.ConfigProvider) gin.HandlerFunc
	JWTAuthMiddleware(token token.Token) gin.HandlerFunc
	RequireRole(role ...string) gin.HandlerFunc
}

type middleware struct {
}

func NewMiddleware() Middleware {
	return &middleware{}
}
