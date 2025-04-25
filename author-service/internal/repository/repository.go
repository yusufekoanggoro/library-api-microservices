package repository

import (
	sharedDomain "author-service/pkg/shared/domain"
	"context"
)

type AuthorRepository interface {
	SaveAuthor(ctx context.Context, author *sharedDomain.Author) error
	GetAllAuthors(ctx context.Context, page, limit int) ([]*sharedDomain.Author, int64, error)
	GetAuthorByID(ctx context.Context, id uint) (*sharedDomain.Author, error)
	DeleteAuthor(ctx context.Context, id uint) error
}
