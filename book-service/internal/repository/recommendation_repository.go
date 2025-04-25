package repository

import (
	"book-service/pkg/shared/domain"
	"context"

	"gorm.io/gorm"
)

type RecommendationRepository interface {
	Create(ctx context.Context, recommendation *domain.Recommendation) error
	GetByID(ctx context.Context, id uint) (*domain.Recommendation, error)
	GetAll(ctx context.Context) ([]*domain.Recommendation, error)
	Delete(ctx context.Context, id uint) error
	FindByUserAndBook(ctx context.Context, userID, bookID uint) (*domain.Recommendation, error)
}

type recommendationRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewRecommendationRepository(readDB *gorm.DB, writeDB *gorm.DB) RecommendationRepository {
	return &recommendationRepository{readDB: readDB, writeDB: writeDB}
}

func (r *recommendationRepository) Create(ctx context.Context, recommendation *domain.Recommendation) error {
	return r.writeDB.Create(recommendation).Error
}

func (r *recommendationRepository) GetByID(ctx context.Context, id uint) (*domain.Recommendation, error) {
	var recommendation domain.Recommendation
	if err := r.readDB.Preload("Book").Preload("User").First(&recommendation, id).Error; err != nil {
		return nil, err
	}
	return &recommendation, nil
}

func (r *recommendationRepository) GetAll(ctx context.Context) ([]*domain.Recommendation, error) {
	var recommendations []*domain.Recommendation
	if err := r.readDB.Find(&recommendations).Error; err != nil {
		return nil, err
	}
	return recommendations, nil
}

func (r *recommendationRepository) Delete(ctx context.Context, id uint) error {
	return r.writeDB.Delete(&domain.Recommendation{}, id).Error
}

func (r *recommendationRepository) FindByUserAndBook(ctx context.Context, userID, bookID uint) (*domain.Recommendation, error) {
	var recommendation domain.Recommendation

	err := r.readDB.
		Where("user_id = ? AND book_id = ?", userID, bookID).
		First(&recommendation).Error

	if err != nil {
		return nil, err // Error akan di-handle di use case
	}

	return &recommendation, nil
}
