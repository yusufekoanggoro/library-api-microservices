package usecase

import (
	"auth-service/internal/domain"
	sharedDomain "auth-service/pkg/shared/domain"
	"context"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*sharedDomain.User, error)
	GetAllUsers(ctx context.Context, req *domain.PaginationRequest) (*domain.PaginatedResponse, error)
	GetUserByID(ctx context.Context, id uint) (*sharedDomain.User, error)
	UpdateUser(ctx context.Context, req *domain.UpdateUserRequest) (*sharedDomain.User, error)
	DeleteUser(ctx context.Context, id uint) error
}

type AuthUsecase interface {
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenResponse, error)
}
