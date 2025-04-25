package usecase

import (
	"author-service/internal/domain"
	sharedDomain "author-service/pkg/shared/domain"
	"context"
)

type AuthorUsecase interface {
	CreateAuthor(ctx context.Context, req *domain.CreateAuthorRequest) (*sharedDomain.Author, error)
	GetAllAuthors(ctx context.Context, req *domain.PaginationRequest) (*domain.PaginatedResponse, error)
	GetAuthorByID(ctx context.Context, id uint) (*sharedDomain.Author, error)
	UpdateAuthor(ctx context.Context, req *domain.UpdateAuthorRequest) (*sharedDomain.Author, error)
	DeleteAuthor(ctx context.Context, id uint) error
}
