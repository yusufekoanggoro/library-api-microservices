package repository

import (
	"book-service/pkg/shared/domain"
	"context"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(ctx context.Context, book *domain.Book) error
	FindByID(ctx context.Context, id uint) (*domain.Book, error)
	FindAll(ctx context.Context, page, limit int, search string) ([]*domain.Book, int64, error)
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id uint) error
}

type bookRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewBookRepository(readDB *gorm.DB, writeDB *gorm.DB) BookRepository {
	return &bookRepository{readDB: readDB, writeDB: writeDB}
}

func (r *bookRepository) Create(ctx context.Context, book *domain.Book) error {
	return r.writeDB.Create(book).Error
}

func (r *bookRepository) FindByID(ctx context.Context, id uint) (*domain.Book, error) {
	var book domain.Book
	if err := r.readDB.
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("Authors").
		Preload("Categories").
		Preload("Categories").
		Preload("Borrowings").
		First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) FindAll(ctx context.Context, page, limit int, search string) ([]*domain.Book, int64, error) {
	var books []*domain.Book
	var total int64

	query := r.readDB.Model(&domain.Book{})
	if search != "" {
		query = query.Where("title ILIKE ? OR isbn ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit

	if err := query.Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		return nil, 0, err
	}
	return books, total, nil
}

func (r *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	return r.writeDB.Save(book).Error
}

func (r *bookRepository) Delete(ctx context.Context, id uint) error {
	return r.writeDB.Delete(&domain.Book{}, id).Error
}
