package usecase

import (
	"book-service/internal/domain"
	sharedDomain "book-service/pkg/shared/domain"
	"book-service/proto/book"
	"context"
)

type BookUsecase interface {
	CreateBook(ctx context.Context, req *domain.CreateBookRequest) (*sharedDomain.Book, error)
	GetBookByID(ctx context.Context, id uint) (*sharedDomain.Book, error)
	ListBooks(ctx context.Context, req *domain.PaginationRequest) (*domain.PaginatedResponse, error)
	UpdateBook(ctx context.Context, req *domain.UpdateBookRequest) (*sharedDomain.Book, error)
	DeleteBook(ctx context.Context, id uint) error
}

type StockUsecase interface {
	IncreaseStock(ctx context.Context, req *domain.IncreaseStockRequest) (*sharedDomain.Book, error)
	DecreaseStock(ctx context.Context, req *domain.DecreaseStockRequest) (*sharedDomain.Book, error)
}

type BookCategoryUsecase interface {
	DeleteBookCategory(categoryID uint) error
}

type BookAuthorUsecase interface {
	DeleteBookAuthor(authorID uint) error
}

type CategoryUsecase interface {
	SaveCategory(ctx context.Context, req *book.CategoryData) error
	DeleteCategory(ctx context.Context, categoryID uint) error
}

type AuthorUsecase interface {
	SaveAuthor(ctx context.Context, req *book.AuthorData) error
	DeleteAuthor(ctx context.Context, authorID uint) error
}

type UserUsecase interface {
	SaveUser(ctx context.Context, req *book.UserData) error
	DeleteUser(ctx context.Context, userID uint) error
}

type BorrowingUseCase interface {
	BorrowBook(ctx context.Context, req *domain.BorrowBookRequest) (*sharedDomain.Borrowing, error)
	ReturnBook(ctx context.Context, req *domain.ReturnBookRequest) (*sharedDomain.Borrowing, error)
	ListBorrowings(ctx context.Context) ([]*sharedDomain.Borrowing, error)
}

type RecommendationUsecase interface {
	CreateRecommendation(ctx context.Context, req *domain.CreateRecommendationRequest) (*sharedDomain.Recommendation, error)
	GetRecommendationByID(ctx context.Context, id uint) (*sharedDomain.Recommendation, error)
	GetAllRecommendations(ctx context.Context) ([]*sharedDomain.Recommendation, error)
	DeleteRecommendation(ctx context.Context, id uint) error
}
