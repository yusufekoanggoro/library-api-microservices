package repository

import (
	"book-service/pkg/shared/domain"
	"context"

	"gorm.io/gorm"
)

type BorrowingRepository interface {
	Create(ctx context.Context, borrowing *domain.Borrowing) error
	FindByID(ctx context.Context, id uint) (*domain.Borrowing, error)
	FindAll(ctx context.Context) ([]*domain.Borrowing, error)
	Update(ctx context.Context, borrowing *domain.Borrowing) error
	Delete(ctx context.Context, id uint) error
	FindByUserAndBook(ctx context.Context, userID, bookID uint) (*domain.Borrowing, error)
	FindByUserAndStatus(ctx context.Context, userID uint, status string) (*domain.Borrowing, error)
}

type borrowingRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewBorrowingRepository(readDB *gorm.DB, writeDB *gorm.DB) BorrowingRepository {
	return &borrowingRepository{readDB: readDB, writeDB: writeDB}
}

func (r *borrowingRepository) Create(ctx context.Context, borrowing *domain.Borrowing) error {
	return r.writeDB.Create(borrowing).Error
}

func (r *borrowingRepository) FindByID(ctx context.Context, id uint) (*domain.Borrowing, error) {
	var borrowing domain.Borrowing
	if err := r.readDB.Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("User").
		Preload("Book").
		First(&borrowing, id).Error; err != nil {
		return nil, err
	}
	return &borrowing, nil
}

func (r *borrowingRepository) FindAll(ctx context.Context) ([]*domain.Borrowing, error) {
	var borrowings []*domain.Borrowing
	if err := r.readDB.Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("User").
		Preload("Book").
		Find(&borrowings).Error; err != nil {
		return nil, err
	}
	return borrowings, nil
}

func (r *borrowingRepository) Update(ctx context.Context, borrowing *domain.Borrowing) error {
	return r.writeDB.Save(borrowing).Error
}

func (r *borrowingRepository) Delete(ctx context.Context, id uint) error {
	return r.writeDB.Delete(&domain.Borrowing{}, id).Error
}

func (r *borrowingRepository) FindByUserAndBook(ctx context.Context, userID, bookID uint) (*domain.Borrowing, error) {
	var borrowing domain.Borrowing

	err := r.readDB.
		Where("user_id = ? AND book_id = ?", userID, bookID).
		First(&borrowing).Error

	if err != nil {
		return nil, err // Error akan di-handle di use case
	}

	return &borrowing, nil
}

func (r *borrowingRepository) FindByUserAndStatus(ctx context.Context, userID uint, status string) (*domain.Borrowing, error) {
	var borrowing domain.Borrowing

	err := r.readDB.
		Where("user_id = ? AND status = ?", userID, status).
		First(&borrowing).Error

	if err != nil {
		return nil, err // Error akan di-handle di use case
	}

	return &borrowing, nil
}
