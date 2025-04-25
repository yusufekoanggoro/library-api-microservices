package domain

type BookCategory struct {
	BookID     uint `gorm:"primaryKey" json:"bookId"`
	CategoryID uint `gorm:"primaryKey" json:"categoryId"`
}
