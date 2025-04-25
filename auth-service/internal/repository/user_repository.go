package repository

import (
	sharedDomain "auth-service/pkg/shared/domain"
	"context"
	"log"

	"gorm.io/gorm"
)

type userRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func NewUserRepository(readDB *gorm.DB, writeDB *gorm.DB) UserRepository {
	return &userRepository{readDB, writeDB}
}

func (r *userRepository) GetAllUsers(ctx context.Context, page, limit int) ([]*sharedDomain.User, int64, error) {
	var users []*sharedDomain.User
	var totalRows int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	if err := r.readDB.Model(&sharedDomain.User{}).Count(&totalRows).Error; err != nil {
		log.Println("GetAllUsers count error:", err)
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := r.readDB.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error
	if err != nil {
		log.Println("GetAllUsers query error:", err)
		return nil, 0, err
	}

	return users, totalRows, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*sharedDomain.User, error) {
	var user sharedDomain.User

	err := r.readDB.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) SaveUser(ctx context.Context, user *sharedDomain.User) error {
	return r.writeDB.Save(user).Error
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint) error {
	return r.writeDB.Where("id = ? ", id).Delete(&sharedDomain.User{}).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*sharedDomain.User, error) {
	var user sharedDomain.User
	if err := r.readDB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
