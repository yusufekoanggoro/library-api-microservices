package domain

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"not null" json:"title"`
	ISBN        string `gorm:"unique;not null" json:"isbn"`
	PublishYear int    `gorm:"not null" json:"publishYear"`
	Stock       int    `gorm:"not null;default:0" json:"stock"`

	CreatedByID uint  `gorm:"not null" json:"createdById"`
	CreatedBy   *User `gorm:"foreignKey:CreatedByID" json:"createdBy"`

	UpdatedByID uint  `gorm:"not null" json:"updatedById"`
	UpdatedBy   *User `gorm:"foreignKey:UpdatedByID" json:"updatedBy"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Authors         []*Author         `gorm:"many2many:book_authors;constraint:OnDelete:CASCADE;" json:"authors"`
	Categories      []*Category       `gorm:"many2many:book_categories;constraint:OnDelete:CASCADE;" json:"categories"`
	Borrowings      []*Borrowing      `gorm:"foreignKey:BookID" json:"borrowings"`
	Recommendations []*Recommendation `gorm:"foreignKey:BookID" json:"recommendations"`
}
