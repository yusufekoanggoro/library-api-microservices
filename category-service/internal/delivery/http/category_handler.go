package http

import (
	"category-service/internal/domain"
	"category-service/internal/usecase"
	"category-service/pkg/shared/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	usecase usecase.CategoryUsecase
}

func NewCategoryHandler(uc usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{usecase: uc}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
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

	var req domain.CreateCategoryRequest
	req.CreatedById = userIDUint
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.usecase.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create category")
		return
	}

	response.Success(c, http.StatusCreated, "Category created successfully", book)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	var req domain.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	categories, err := h.usecase.GetAllCategories(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to retrieve categories")
		return
	}

	pagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    req.Limit,
		TotalPages:  categories.TotalPages,
		TotalItems:  int(categories.Total),
	}

	response.SuccessWithPagination(c, http.StatusOK, "Categories retrieved successfully", categories.Data, pagination)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	category, err := h.usecase.GetCategoryByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Category not found")
		return
	}

	response.Success(c, http.StatusOK, "Category retrieved successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
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

	var req domain.UpdateCategoryRequest
	req.ID = uint(categoryID)
	req.UpdatedById = userIDUint
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.usecase.UpdateCategory(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update category")
		return
	}

	response.Success(c, http.StatusOK, "Category updated successfully", book)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.usecase.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	response.Success(c, http.StatusOK, "Category deleted successfully", nil)
}
