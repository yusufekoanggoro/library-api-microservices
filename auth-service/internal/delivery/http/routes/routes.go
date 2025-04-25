package routes

import (
	"auth-service/config"
	"auth-service/internal/delivery/http"
	"auth-service/pkg/logger"
	"auth-service/pkg/middleware"
	"auth-service/pkg/token"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	cfg config.ConfigProvider,
	authHandler http.AuthHandler,
	userHandler http.UserHandler,
	logger logger.Logger,
	token token.Token,
	middleware middleware.Middleware,
) *gin.Engine {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", middleware.BasicAuthMiddleware(cfg), authHandler.Login)
		authGroup.POST("/refresh-token", middleware.BasicAuthMiddleware(cfg), authHandler.RefreshToken)
	}

	userRoutes := r.Group("/users", middleware.JWTAuthMiddleware(token), middleware.RequireRole("admin"))
	{
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/:id", userHandler.GetUserByID)
		userRoutes.PATCH("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	return r
}
