package repository

import (
	"book-service/pkg/shared/domain"
	"context"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Upsert(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uint) error
}

type categoryRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewCategoryRepository(readDB *gorm.DB, writeDB *gorm.DB) CategoryRepository {
	return &categoryRepository{readDB: readDB, writeDB: writeDB}
}

func (r *categoryRepository) Upsert(ctx context.Context, category *domain.Category) error {
	return r.writeDB.Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	return r.writeDB.Delete(&domain.Category{}, id).Error
}
