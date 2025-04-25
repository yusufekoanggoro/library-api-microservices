package token

import (
	"category-service/internal/domain"
	"time"
)

type Token interface {
	GenerateToken(userId uint, role string, expired time.Duration) (string, error)
	ValidateToken(tokenString string) (*domain.TokenClaims, error)
}
