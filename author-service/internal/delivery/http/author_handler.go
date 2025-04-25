package http

import (
	"author-service/internal/domain"
	"author-service/internal/usecase"
	"author-service/pkg/shared/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthorHandler struct {
	usecase usecase.AuthorUsecase
}

func NewAuthorHandler(uc usecase.AuthorUsecase) *AuthorHandler {
	return &AuthorHandler{usecase: uc}
}

func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var userID interface{}
	var exists bool

	if userID, exists = c.Get("userId"); !exists {
		response.Error(c, http.StatusBadRequest, "User ID not found")
		return
	}

	userIdUint, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid User ID type")
		return
	}

	var req domain.CreateAuthorRequest
	req.CreatedByID = userIdUint
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.usecase.CreateAuthor(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create author")
		return
	}

	response.Success(c, http.StatusCreated, "Creator created successfully", book)
}

func (h *AuthorHandler) GetAllAuthors(c *gin.Context) {
	var req domain.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	authors, err := h.usecase.GetAllAuthors(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to retrieve authors")
		return
	}

	pagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    req.Limit,
		TotalPages:  authors.TotalPages,
		TotalItems:  int(authors.Total),
	}

	response.SuccessWithPagination(c, http.StatusOK, "Authors retrieved successfully", authors.Data, pagination)
}

func (h *AuthorHandler) GetAuthorByID(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	author, err := h.usecase.GetAuthorByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Author not found")
		return
	}

	response.Success(c, http.StatusOK, "Author retrieved successfully", author)
}

func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	authorID, err := strconv.Atoi(c.Param("id"))
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
		response.Error(c, http.StatusInternalServerError, "Invalid User ID type")
		return
	}

	var req domain.UpdateAuthorRequest
	req.ID = uint(authorID)
	req.UpdatedByID = userIDUint
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.usecase.UpdateAuthor(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update author")
		return
	}

	response.Success(c, http.StatusOK, "Author updated successfully", book)
}

func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.usecase.DeleteAuthor(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete author")
		return
	}

	response.Success(c, http.StatusOK, "Author deleted successfully", nil)
}
