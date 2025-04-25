package usecase

import (
	"book-service/internal/domain"
	"book-service/internal/repository"
	sharedDomain "book-service/pkg/shared/domain"
	"context"
)

type bookUsecase struct {
	repo repository.Repository
}

func NewBookUsecase(
	repo repository.Repository,
) BookUsecase {
	return &bookUsecase{
		repo: repo,
	}
}

func (uc *bookUsecase) CreateBook(ctx context.Context, req *domain.CreateBookRequest) (*sharedDomain.Book, error) {
	book := &sharedDomain.Book{
		Title:       req.Title,
		ISBN:        req.ISBN,
		PublishYear: req.PublishYear,
		Stock:       req.Stock,
		CreatedByID: req.CreatedByID,
		UpdatedByID: req.CreatedByID,
	}

	err := uc.repo.WithTransaction(func(txRepo repository.Repository) error {
		err := txRepo.Book().Create(ctx, book)
		if err != nil {
			return err
		}

		var bookAuthors []*sharedDomain.BookAuthor
		for _, authorID := range req.AuthorIDs {
			bookAuthors = append(bookAuthors, &sharedDomain.BookAuthor{BookID: book.ID, AuthorID: authorID})
		}

		if len(bookAuthors) > 0 {
			if err := txRepo.BookAuthor().Create(ctx, bookAuthors); err != nil {
				return err
			}
		}

		var bookCategories []*sharedDomain.BookCategory
		for _, categoryID := range req.CategoryIDs {
			bookCategories = append(bookCategories, &sharedDomain.BookCategory{BookID: book.ID, CategoryID: categoryID})
		}

		if len(bookCategories) > 0 {
			if err := txRepo.BookCategory().Create(ctx, bookCategories); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (uc *bookUsecase) GetBookByID(ctx context.Context, id uint) (*sharedDomain.Book, error) {
	return uc.repo.Book().FindByID(ctx, id)
}

func (uc *bookUsecase) ListBooks(ctx context.Context, req *domain.PaginationRequest) (*domain.PaginatedResponse, error) {
	books, totalRows, err := uc.repo.Book().FindAll(ctx, req.Page, req.Limit, req.Search)
	if err != nil {
		return nil, err
	}

	paginatedResponse := &domain.PaginatedResponse{
		Data:       books,
		Total:      totalRows,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((totalRows + int64(req.Limit) - 1) / int64(req.Limit)),
	}

	return paginatedResponse, nil
}

func (uc *bookUsecase) UpdateBook(ctx context.Context, req *domain.UpdateBookRequest) (*sharedDomain.Book, error) {
	book, err := uc.repo.Book().FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.ISBN != nil {
		book.ISBN = *req.ISBN
	}
	if req.PublishYear != nil {
		book.PublishYear = *req.PublishYear
	}
	if req.Stock != nil {
		book.Stock = *req.Stock
	}

	err = uc.repo.WithTransaction(func(txRepo repository.Repository) error {
		if err := txRepo.Book().Update(ctx, book); err != nil {
			return err
		}

		if req.AuthorIDs != nil {
			if err := txRepo.BookAuthor().DeleteByBookID(ctx, book.ID); err != nil {
				return err
			}

			var bookAuthors []*sharedDomain.BookAuthor
			for _, authorID := range *req.AuthorIDs {
				bookAuthors = append(bookAuthors, &sharedDomain.BookAuthor{BookID: book.ID, AuthorID: authorID})
			}
			if len(bookAuthors) > 0 {
				if err := txRepo.BookAuthor().Create(ctx, bookAuthors); err != nil {
					return err
				}
			}
		}

		if req.CategoryIDs != nil {
			if err := txRepo.BookCategory().DeleteByBookID(ctx, book.ID); err != nil {
				return err
			}

			var bookCategories []*sharedDomain.BookCategory
			for _, categoryID := range *req.CategoryIDs {
				bookCategories = append(bookCategories, &sharedDomain.BookCategory{BookID: book.ID, CategoryID: categoryID})
			}

			if len(bookCategories) > 0 {
				if err := txRepo.BookCategory().Create(ctx, bookCategories); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return uc.repo.Book().FindByID(ctx, req.ID)

}

func (uc *bookUsecase) DeleteBook(ctx context.Context, id uint) error {
	return uc.repo.Book().Delete(ctx, id)
}
