package repository

import (
	"book-service/pkg/shared/domain"
	"context"

	"gorm.io/gorm"
)

type BookCategoryRepository interface {
	Create(ctx context.Context, bookCategory []*domain.BookCategory) error
	DeleteByBookID(ctx context.Context, bookID uint) error
}

type bookCategoryRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewBookCategoryRepository(readDB *gorm.DB, writeDB *gorm.DB) BookCategoryRepository {
	return &bookCategoryRepository{readDB: readDB, writeDB: writeDB}
}

func (r *bookCategoryRepository) Create(ctx context.Context, bookCategory []*domain.BookCategory) error {
	return r.writeDB.Create(bookCategory).Error
}

func (r *bookCategoryRepository) DeleteByBookID(ctx context.Context, bookID uint) error {
	return r.writeDB.Where("book_id = ?", bookID).Delete(&domain.BookCategory{}).Error
}
