package repository

import (
	sharedDomain "auth-service/pkg/shared/domain"
	"context"
)

type UserRepository interface {
	SaveUser(ctx context.Context, category *sharedDomain.User) error
	GetAllUsers(ctx context.Context, page, limit int) ([]*sharedDomain.User, int64, error)
	GetUserByID(ctx context.Context, id uint) (*sharedDomain.User, error)
	DeleteUser(ctx context.Context, id uint) error
	FindByEmail(ctx context.Context, email string) (*sharedDomain.User, error)
}
