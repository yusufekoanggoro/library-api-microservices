package http

import (
	"book-service/internal/domain"
	"book-service/internal/usecase"
	"net/http"
	"strconv"

	"book-service/pkg/shared/response"

	"github.com/gin-gonic/gin"
)

type BookHandler interface {
	CreateBook(c *gin.Context)
	GetListBook(c *gin.Context)
	GetBookByID(c *gin.Context)
	UpdateBook(c *gin.Context)
	DeleteBook(c *gin.Context)
}

type bookHandler struct {
	usecase usecase.BookUsecase
}

func NewBookHandler(uc usecase.BookUsecase) BookHandler {
	return &bookHandler{usecase: uc}
}

func (h *bookHandler) CreateBook(c *gin.Context) {
	var userID interface{}
	var exists bool

	if userID, exists = c.Get("userId"); !exists {
		response.Error(c, http.StatusBadRequest, "User ID not found")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid User ID type")
		return
	}

	var req domain.CreateBookRequest
	req.CreatedByID = userIDUint

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.usecase.CreateBook(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create book")
		return
	}

	response.Success(c, http.StatusCreated, "Book created successfully", book)
}

func (h *bookHandler) GetListBook(c *gin.Context) {
	var req domain.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	books, err := h.usecase.ListBooks(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to retrieve books")
		return
	}

	pagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    req.Limit,
		TotalPages:  books.TotalPages,
		TotalItems:  int(books.Total),
	}

	response.SuccessWithPagination(c, http.StatusOK, "Books retrieved successfully", books.Data, pagination)
}

func (h *bookHandler) GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.usecase.GetBookByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Book not found")
		return
	}

	response.Success(c, http.StatusOK, "Book retrieved successfully", book)
}

func (h *bookHandler) UpdateBook(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var userID interface{}
	var exists bool

	if userID, exists = c.Get("userId"); !exists {
		response.Error(c, http.StatusBadRequest, "User ID not found")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid User ID type")
		return
	}

	var req domain.UpdateBookRequest
	req.ID = uint(bookID)
	req.UpdatedByID = userIDUint
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.usecase.UpdateBook(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update book")
		return
	}

	response.Success(c, http.StatusOK, "Book updated successfully", book)
}

func (h *bookHandler) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.usecase.DeleteBook(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete book")
		return
	}

	response.Success(c, http.StatusOK, "Book deleted successfully", nil)
}
