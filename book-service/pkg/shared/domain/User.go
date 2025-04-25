package domain

import (
	"time"

	protoBook "book-service/proto/book"

	"gorm.io/gorm"
)

type User struct {
	ID              uint              `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Username        string            `gorm:"unique;not null" json:"username"`
	Email           string            `gorm:"unique;not null" json:"email"`
	Role            string            `gorm:"not null" json:"role"`
	Password        string            `gorm:"not null" json:"-"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"deletedAt,omitempty"`
	Borrowings      []*Borrowing      `gorm:"foreignKey:BookID" json:"borrowings"`
	Recommendations []*Recommendation `gorm:"foreignKey:BookID" json:"recommendations"`
}

func (m *User) ToProto() *protoBook.UserData {
	var deletedAt string
	if m.DeletedAt.Valid {
		deletedAt = m.DeletedAt.Time.Format(time.RFC3339)
	}

	return &protoBook.UserData{
		Id:       uint32(m.ID),
		Username: m.Username,
		Email:    m.Email,
		Role:     m.Role,
		Password: m.Password,

		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
		DeletedAt: deletedAt,
	}
}

func (m *User) FromProto(pb *protoBook.UserData) (*User, error) {
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

	return &User{
		ID:       uint(pb.Id),
		Username: pb.Username,
		Email:    pb.Email,
		Role:     pb.Role,
		Password: pb.Password,

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}, nil
}
