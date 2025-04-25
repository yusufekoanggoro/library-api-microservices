package repository

import (
	"book-service/pkg/shared/domain"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	Upsert(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*domain.User, error)
}

type userRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewUserRepository(readDB *gorm.DB, writeDB *gorm.DB) UserRepository {
	return &userRepository{readDB: readDB, writeDB: writeDB}
}

func (r *userRepository) Upsert(ctx context.Context, user *domain.User) error {
	return r.writeDB.Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.writeDB.Delete(&domain.User{}, id).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.readDB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
