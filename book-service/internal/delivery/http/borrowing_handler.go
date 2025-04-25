package http

import (
	"book-service/internal/domain"
	"book-service/internal/usecase"
	"book-service/pkg/shared/response"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BorrowingHandler interface {
	BorrowBook(c *gin.Context)
	ReturnBook(c *gin.Context)
	GetListBorrowing(c *gin.Context)
}

type borrowingHandler struct {
	useCase usecase.BorrowingUseCase
}

func NewBorrowingHandler(useCase usecase.BorrowingUseCase) BorrowingHandler {
	return &borrowingHandler{useCase: useCase}
}

func (h *borrowingHandler) BorrowBook(c *gin.Context) {
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

	var req domain.BorrowBookRequest
	req.CreatedByID = userIDUint

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ctx := context.Background()
	borrowing, err := h.useCase.BorrowBook(ctx, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to borrow book")
		return
	}

	response.Success(c, http.StatusOK, "Book successfully borrowed", borrowing)
}

func (h *borrowingHandler) ReturnBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
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

	var req domain.ReturnBookRequest
	req.ID = uint(id)
	req.UpdatedByID = userIDUint
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx := context.Background()
	borrowing, err := h.useCase.ReturnBook(ctx, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to return book")
		return
	}

	response.Success(c, http.StatusOK, "Book successfully returned", borrowing)
}

func (h *borrowingHandler) GetListBorrowing(c *gin.Context) {
	borrowings, err := h.useCase.ListBorrowings(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to retrieve borrowings")
		return
	}

	response.Success(c, http.StatusOK, "Borrowings retrieved successfully", borrowings)
}
