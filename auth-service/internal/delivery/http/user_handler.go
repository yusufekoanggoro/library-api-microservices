package http

import (
	"auth-service/internal/domain"
	"auth-service/internal/usecase"
	"auth-service/pkg/shared/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userHandler struct {
	usecase usecase.UserUsecase
}

func NeUserHandler(uc usecase.UserUsecase) UserHandler {
	return &userHandler{usecase: uc}
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.usecase.CreateUser(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", user)
}

func (h *userHandler) GetAllUsers(c *gin.Context) {
	var req domain.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	users, err := h.usecase.GetAllUsers(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to retrieve users")
		return
	}

	pagination := response.Pagination{
		CurrentPage: req.Page,
		PageSize:    req.Limit,
		TotalPages:  users.TotalPages,
		TotalItems:  int(users.Total),
	}

	response.SuccessWithPagination(c, http.StatusOK, "Users retrieved successfully", users.Data, pagination)
}

func (h *userHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.usecase.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

func (h *userHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var req domain.UpdateUserRequest
	req.ID = uint(userID)
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := h.usecase.UpdateUser(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", book)
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.usecase.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
