package http

import (
	"auth-service/internal/domain"
	"auth-service/internal/usecase"
	"auth-service/pkg/shared/response"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type authHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(usecase usecase.AuthUsecase) AuthHandler {
	return &authHandler{usecase: usecase}
}

func (h *authHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := h.usecase.Login(context.Background(), &req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	response.Success(c, http.StatusOK, "Login successful", resp)
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := h.usecase.RefreshToken(context.Background(), req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", resp)
}
