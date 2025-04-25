package seeders

import (
	sharedDomain "auth-service/pkg/shared/domain"

	"gorm.io/gorm"
)

func SeedUsers(readDB *gorm.DB, writeDB *gorm.DB) ([]*sharedDomain.User, error) {
	users := []*sharedDomain.User{
		{Username: "admin", Email: "admin@admin.com", Password: "admin123", Role: "admin"},
	}

	for _, user := range users {
		if err := user.HashPassword(); err != nil {
			return nil, err
		}

		if err := writeDB.Create(&user).Error; err != nil {
			return nil, err
		}
	}

	return users, nil
}
