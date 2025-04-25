package repository

import (
	"gorm.io/gorm"
)

type Repository interface {
	Book() BookRepository
	Stock() StockRepository
	Category() CategoryRepository
	Author() AuthorRepository
	User() UserRepository
	Borrowing() BorrowingRepository
	BookAuthor() BookAuthorRepository
	BookCategory() BookCategoryRepository
	Recommendation() RecommendationRepository
	WithTransaction(fn func(txRepo Repository) error) error
}

type repository struct {
	writeDB        *gorm.DB
	book           BookRepository
	stock          StockRepository
	category       CategoryRepository
	author         AuthorRepository
	user           UserRepository
	borrowing      BorrowingRepository
	bookAuthor     BookAuthorRepository
	bookCategory   BookCategoryRepository
	recommendation RecommendationRepository
}

func NewRepository(readDB *gorm.DB, writeDB *gorm.DB) Repository {
	return &repository{
		writeDB:        writeDB,
		book:           NewBookRepository(readDB, writeDB),
		stock:          NewStockRepository(readDB, writeDB),
		category:       NewCategoryRepository(readDB, writeDB),
		author:         NewAuthorRepository(readDB, writeDB),
		user:           NewUserRepository(readDB, writeDB),
		borrowing:      NewBorrowingRepository(readDB, writeDB),
		bookAuthor:     NewBookAuthorRepository(readDB, writeDB),
		bookCategory:   NewBookCategoryRepository(readDB, writeDB),
		recommendation: NewRecommendationRepository(readDB, writeDB),
	}
}

func (r *repository) Book() BookRepository                     { return r.book }
func (r *repository) Stock() StockRepository                   { return r.stock }
func (r *repository) Category() CategoryRepository             { return r.category }
func (r *repository) Author() AuthorRepository                 { return r.author }
func (r *repository) User() UserRepository                     { return r.user }
func (r *repository) Borrowing() BorrowingRepository           { return r.borrowing }
func (r *repository) BookAuthor() BookAuthorRepository         { return r.bookAuthor }
func (r *repository) BookCategory() BookCategoryRepository     { return r.bookCategory }
func (r *repository) Recommendation() RecommendationRepository { return r.recommendation }

func (r *repository) WithTransaction(fn func(txRepo Repository) error) error {
	return r.writeDB.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx, tx)
		return fn(txRepo)
	})
}
