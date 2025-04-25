package repository

import (
	sharedDomain "category-service/pkg/shared/domain"
	"context"
	"log"

	"gorm.io/gorm"
)

type categoryRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewAuthorRepository(readDB *gorm.DB, writeDB *gorm.DB) CategoryRepository {
	return &categoryRepository{readDB: readDB, writeDB: writeDB}
}

func (r *categoryRepository) GetAllCategories(ctx context.Context, page, limit int) ([]*sharedDomain.Category, int64, error) {
	var categories []*sharedDomain.Category
	var totalRows int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	if err := r.readDB.Model(&sharedDomain.Category{}).Count(&totalRows).Error; err != nil {
		log.Println("GetAllCategories count error:", err)
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := r.readDB.Limit(limit).Offset(offset).Order("created_at DESC").Find(&categories).Error
	if err != nil {
		log.Println("GetAllCategories query error:", err)
		return nil, 0, err
	}

	return categories, totalRows, nil
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id uint) (*sharedDomain.Category, error) {
	var category sharedDomain.Category

	err := r.readDB.First(&category, id).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) SaveCategory(ctx context.Context, category *sharedDomain.Category) error {
	return r.writeDB.Save(category).Error
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id uint) error {
	return r.writeDB.Where("id = ? ", id).Delete(&sharedDomain.Category{}).Error
}
