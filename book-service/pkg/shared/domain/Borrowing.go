package domain

import (
	"time"

	"gorm.io/gorm"
)

type Borrowing struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null" json:"userId"`
	BookID     uint       `gorm:"not null" json:"bookId"`
	BorrowDate time.Time  `gorm:"not null" json:"borrowDate"`
	ReturnDate *time.Time `gorm:"default:null" json:"returnDate"`

	Status string `gorm:"type:varchar(20);not null;default:'borrowed'" json:"status"`

	CreatedByID uint  `gorm:"not null" json:"createdById"`
	CreatedBy   *User `gorm:"foreignKey:CreatedByID" json:"createdBy"`

	UpdatedByID uint  `gorm:"not null" json:"updatedById"`
	UpdatedBy   *User `gorm:"foreignKey:UpdatedByID" json:"updatedBy"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"user"`
	Book *Book `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"book"`
}
