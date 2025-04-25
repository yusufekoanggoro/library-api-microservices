package domain

import (
	"time"

	"gorm.io/gorm"
)

type Recommendation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null" json:"userId"`
	BookID    uint           `gorm:"not null" json:"bookId"`
	Reason    string         `gorm:"not null" json:"reason"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	User *User `gorm:"foreignKey:UserID" json:"user"`
	Book *Book `gorm:"foreignKey:BookID" json:"book"`
}
