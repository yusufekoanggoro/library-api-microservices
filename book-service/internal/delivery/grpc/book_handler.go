package grpc

import (
	"book-service/internal/usecase"
	"book-service/proto/book"
	"context"
)

type BookHandler struct {
	book.UnimplementedBookServiceServer
	userUsecase     usecase.UserUsecase
	authorUsecase   usecase.AuthorUsecase
	categoryUsecase usecase.CategoryUsecase
}

func NewBookHandler(userUsecase usecase.UserUsecase, authorUsecase usecase.AuthorUsecase, categoryUsecase usecase.CategoryUsecase) *BookHandler {
	return &BookHandler{userUsecase: userUsecase, authorUsecase: authorUsecase, categoryUsecase: categoryUsecase}
}

func (h *BookHandler) ReceiveUser(c context.Context, req *book.UserData) (*book.BookResponse, error) {
	err := h.userUsecase.SaveUser(c, req)
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &book.BookResponse{
		Success: true,
		Message: "User saved successfully",
	}, nil
}

func (h *BookHandler) DeleteUser(c context.Context, req *book.DeleteData) (*book.BookResponse, error) {
	err := h.userUsecase.DeleteUser(c, uint(req.Id))
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &book.BookResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}

func (h *BookHandler) ReceiveAuthor(c context.Context, req *book.AuthorData) (*book.BookResponse, error) {
	err := h.authorUsecase.SaveAuthor(c, req)
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &book.BookResponse{
		Success: true,
		Message: "Author saved successfully",
	}, nil
}

func (h *BookHandler) DeleteAuthor(c context.Context, req *book.DeleteData) (*book.BookResponse, error) {
	err := h.authorUsecase.DeleteAuthor(c, uint(req.Id))
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &book.BookResponse{
		Success: true,
		Message: "Author deleted successfully",
	}, nil
}

func (h *BookHandler) ReceiveCategory(c context.Context, req *book.CategoryData) (*book.BookResponse, error) {
	err := h.categoryUsecase.SaveCategory(c, req)
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &book.BookResponse{
		Success: true,
		Message: "Category saved successfully",
	}, nil
}

func (h *BookHandler) DeleteCategory(c context.Context, req *book.DeleteData) (*book.BookResponse, error) {
	err := h.categoryUsecase.DeleteCategory(c, uint(req.Id))
	if err != nil {
		return &book.BookResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &book.BookResponse{
		Success: true,
		Message: "Category deleted successfully",
	}, nil
}
