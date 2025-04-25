package http

import (
	"book-service/internal/domain"
	"book-service/internal/usecase"
	"book-service/pkg/shared/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StockHandler interface {
	IncreaseStock(c *gin.Context)
	DecreaseStock(c *gin.Context)
}

type stockHandler struct {
	usecase usecase.StockUsecase
}

func NewStockHandler(uc usecase.StockUsecase) StockHandler {
	return &stockHandler{usecase: uc}
}

func (s *stockHandler) IncreaseStock(c *gin.Context) {
	var req *domain.IncreaseStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := s.usecase.IncreaseStock(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to increase stock")
		return
	}

	response.Success(c, http.StatusOK, "Stock successfully increased", book)
}

func (s *stockHandler) DecreaseStock(c *gin.Context) {
	var req *domain.DecreaseStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book, err := s.usecase.DecreaseStock(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to decrease stock")
		return
	}

	response.Success(c, http.StatusOK, "Stock successfully decreased", book)
}
