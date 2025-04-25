package repository

import (
	"book-service/pkg/shared/domain"
	"context"

	"gorm.io/gorm"
)

type AuthorRepository interface {
	Upsert(ctx context.Context, user *domain.Author) error
	Delete(ctx context.Context, id uint) error
}

type authorRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewAuthorRepository(readDB *gorm.DB, writeDB *gorm.DB) AuthorRepository {
	return &authorRepository{readDB: readDB, writeDB: writeDB}
}

func (r *authorRepository) Upsert(ctx context.Context, author *domain.Author) error {
	return r.writeDB.Save(author).Error
}

func (r *authorRepository) Delete(ctx context.Context, id uint) error {
	return r.writeDB.Delete(&domain.Author{}, id).Error
}
