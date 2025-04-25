package http

import (
	"book-service/internal/domain"
	"book-service/internal/usecase"
	"book-service/pkg/shared/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler interface {
	CreateRecommendation(c *gin.Context)
	GetRecommendationByID(c *gin.Context)
	GetAllRecommendations(c *gin.Context)
	DeleteRecommendation(c *gin.Context)
}

type recommendationHandler struct {
	usecase usecase.RecommendationUsecase
}

func NewRecommendationHandler(usecase usecase.RecommendationUsecase) RecommendationHandler {
	return &recommendationHandler{usecase: usecase}
}

func (h *recommendationHandler) CreateRecommendation(c *gin.Context) {
	var req domain.CreateRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	recommendation, err := h.usecase.CreateRecommendation(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create recommendation")
		return
	}

	response.Success(c, http.StatusCreated, "Recommendation created successfully", recommendation)
}

func (h *recommendationHandler) GetRecommendationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	recommendation, err := h.usecase.GetRecommendationByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Recommendation not found")
		return
	}

	response.Success(c, http.StatusOK, "Recommendation retrieved successfully", recommendation)
}

func (h *recommendationHandler) GetAllRecommendations(c *gin.Context) {
	recommendations, err := h.usecase.GetAllRecommendations(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to retrieve recommendations")
		return
	}

	response.Success(c, http.StatusOK, "Recommendations retrieved successfully", recommendations)
}

func (h *recommendationHandler) DeleteRecommendation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.usecase.DeleteRecommendation(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete recommendation")
		return
	}

	response.Success(c, http.StatusOK, "Recommendation deleted successfully", nil)
}
