package domain

type BookAuthor struct {
	BookID   uint `gorm:"primaryKey" json:"bookId"`
	AuthorID uint `gorm:"primaryKey" json:"authorId"`
}
