package domain

import (
	protoBook "book-service/proto/book"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`

	CreatedByID uint  `gorm:"not null" json:"createdById"`
	CreatedBy   *User `gorm:"foreignKey:CreatedByID" json:"createdBy"`

	UpdatedByID uint  `gorm:"not null" json:"updatedById"`
	UpdatedBy   *User `gorm:"foreignKey:UpdatedByID" json:"updatedBy"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (m *Category) ToProto() *protoBook.CategoryData {
	var deletedAt string
	if m.DeletedAt.Valid {
		deletedAt = m.DeletedAt.Time.Format(time.RFC3339)
	}

	return &protoBook.CategoryData{
		Id:   uint32(m.ID),
		Name: m.Name,

		CreatedById: uint32(m.CreatedByID),
		UpdatedById: uint32(m.UpdatedByID),

		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
		DeletedAt: deletedAt,
	}
}

func (m *Category) FromProto(pb *protoBook.CategoryData) (*Category, error) {
	createdAt, err := time.Parse(time.RFC3339, pb.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339, pb.UpdatedAt)
	if err != nil {
		return nil, err
	}

	var deletedAt gorm.DeletedAt
	if pb.DeletedAt != "" {
		t, err := time.Parse(time.RFC3339, pb.DeletedAt)
		if err != nil {
			return nil, err
		}
		deletedAt = gorm.DeletedAt{Time: t, Valid: true}
	}

	return &Category{
		ID:   uint(pb.Id),
		Name: pb.Name,

		CreatedByID: uint(pb.CreatedById),
		UpdatedByID: uint(pb.UpdatedById),

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}, nil
}
