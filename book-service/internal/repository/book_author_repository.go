package repository

import (
	"book-service/pkg/shared/domain"
	"context"

	"gorm.io/gorm"
)

type BookAuthorRepository interface {
	Create(ctx context.Context, bookAuthor []*domain.BookAuthor) error
	DeleteByBookID(ctx context.Context, bookID uint) error
}

type bookAuthorRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewBookAuthorRepository(readDB *gorm.DB, writeDB *gorm.DB) BookAuthorRepository {
	return &bookAuthorRepository{readDB: readDB, writeDB: writeDB}
}

func (r *bookAuthorRepository) Create(ctx context.Context, bookAuthor []*domain.BookAuthor) error {
	return r.writeDB.Create(bookAuthor).Error
}

func (r *bookAuthorRepository) DeleteByBookID(ctx context.Context, bookID uint) error {
	return r.writeDB.Where("book_id = ? ", bookID).Delete(&domain.BookAuthor{}).Error
}
