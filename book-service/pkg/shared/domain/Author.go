package domain

import (
	protoBook "book-service/proto/book"
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Name string `gorm:"not null" json:"name"`
	Bio  string `gorm:"type:text" json:"bio"`

	CreatedByID uint  `gorm:"not null" json:"createdById"`
	CreatedBy   *User `gorm:"foreignKey:CreatedByID" json:"createdBy"`

	UpdatedByID uint  `gorm:"not null" json:"updatedById"`
	UpdatedBy   *User `gorm:"foreignKey:UpdatedByID" json:"updatedBy"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Books []*Book `gorm:"many2many:book_authors;constraint:OnDelete:CASCADE;" json:"books"`
}

func (m *Author) ToProto() *protoBook.AuthorData {
	var deletedAt string
	if m.DeletedAt.Valid {
		deletedAt = m.DeletedAt.Time.Format(time.RFC3339)
	}

	return &protoBook.AuthorData{
		Id:   uint32(m.ID),
		Name: m.Name,
		Bio:  m.Bio,

		CreatedById: uint32(m.CreatedByID),
		UpdatedById: uint32(m.UpdatedByID),

		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
		DeletedAt: deletedAt,
	}
}

func (m *Author) FromProto(pb *protoBook.AuthorData) (*Author, error) {
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

	return &Author{
		ID:   uint(pb.Id),
		Name: pb.Name,
		Bio:  pb.Bio,

		CreatedByID: uint(pb.CreatedById),
		UpdatedByID: uint(pb.UpdatedById),

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}, nil
}
