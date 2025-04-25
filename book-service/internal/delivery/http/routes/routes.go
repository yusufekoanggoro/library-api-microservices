package routes

import (
	"book-service/config"
	"book-service/internal/delivery/http"
	"book-service/pkg/logger"
	"book-service/pkg/middleware"
	"book-service/pkg/token"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	cfg config.ConfigProvider,
	bookHandler http.BookHandler,
	borrowingHandler http.BorrowingHandler,
	recommendationHandler http.RecommendationHandler,
	stockHandler http.StockHandler,
	logger logger.Logger,
	token token.Token,
	middleware middleware.Middleware,
) *gin.Engine {
	bookRoutes := r.Group("/books", middleware.JWTAuthMiddleware(token), middleware.RequireRole("admin"))
	{
		bookRoutes.POST("/", bookHandler.CreateBook)
		bookRoutes.GET("/", bookHandler.GetListBook)
		bookRoutes.GET("/:id", bookHandler.GetBookByID)
		bookRoutes.PATCH("/:id", bookHandler.UpdateBook)
		bookRoutes.DELETE("/:id", bookHandler.DeleteBook)

		bookRoutes.PUT("/increase-stock", stockHandler.IncreaseStock)
		bookRoutes.PUT("/decrease-stock", stockHandler.DecreaseStock)

		bookRoutes.GET("/borrowings", borrowingHandler.GetListBorrowing)
		bookRoutes.POST("/borrowings", borrowingHandler.BorrowBook)
		bookRoutes.PUT("/borrowings/:id/return", borrowingHandler.ReturnBook)
	}

	recommendationRoutes := r.Group("/books/recommendations", middleware.JWTAuthMiddleware(token), middleware.RequireRole("admin", "member"))
	{
		recommendationRoutes.POST("/", recommendationHandler.CreateRecommendation)
		recommendationRoutes.GET("/:id", recommendationHandler.GetRecommendationByID)
		recommendationRoutes.GET("", recommendationHandler.GetAllRecommendations)
		recommendationRoutes.DELETE("/:id", recommendationHandler.DeleteRecommendation)
	}

	return r
}
