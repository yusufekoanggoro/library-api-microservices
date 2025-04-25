package main

import (
	"author-service/config"
	"author-service/config/key"
	deliveryH "author-service/internal/delivery/http"
	"author-service/internal/grpcservice"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"author-service/internal/repository"
	"author-service/internal/usecase"
	"author-service/pkg/database"
	"author-service/pkg/logger"
	"author-service/pkg/middleware"
	sharedDomain "author-service/pkg/shared/domain"
	"author-service/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()

	logger := logger.NewLogger("author-service", logrus.DebugLevel, os.Stdout)

	_, cancel := context.WithTimeout(context.Background(), 1*time.Minute) // timeout to shutdown server and init configuration
	defer func() {
		cancel()
		if r := recover(); r != nil {
			// notify if this service failed to start (to logstash or slack)
			fmt.Println("Failed to start author service:", r)
			fmt.Printf("Stack trace: \n%s\n", debug.Stack())
		}
	}()

	db := &database.GormDatabase{}
	if err := db.Connect(cfg); err != nil {
		logger.Panic(fmt.Sprintf("Database connection error: %v", err), "db-error", "connection")
	}

	if err := db.AutoMigrate(
		&sharedDomain.Author{},
	); err != nil {
		logger.Panic(fmt.Sprintf("Failed to perform migration: %v", err), "migration", "error")
	}

	privateKeyPath := "config/key/private_key.pem"
	publicKeyPath := "config/key/public_key.pem"

	privateKey, publicKey, err := key.LoadRSAKeys(privateKeyPath, publicKeyPath)
	if err != nil {
		logger.Panic(fmt.Sprintf("Failed to load RSA keys: %v", err), "load rsa", "error")
	}

	jwtService := token.NewJWT(publicKey, privateKey)

	bookClient := grpcservice.NewBookGRPCClient(cfg.GetBookGRPCHost() + ":" + cfg.GetBookGRPCPort())

	// Setup repository, usecase, dan handler
	authorRepo := repository.NewAuthorRepository(db.GetDB())
	authorUsecase := usecase.NewAuthorUsecase(authorRepo, bookClient)
	authorHandler := deliveryH.NewAuthorHandler(authorUsecase)

	// Setup routes
	httpServer := gin.Default()

	authorRoutes := httpServer.Group("/authors", middleware.JWTAuthMiddleware(jwtService), middleware.RequireRole("admin"))
	{
		authorRoutes.POST("/", authorHandler.CreateAuthor)
		authorRoutes.GET("/", authorHandler.GetAllAuthors)
		authorRoutes.GET("/:id", authorHandler.GetAuthorByID)
		authorRoutes.PATCH("/:id", authorHandler.UpdateAuthor)
		authorRoutes.DELETE("/:id", authorHandler.DeleteAuthor)
	}

	httpPort := cfg.GetHTTPPort()
	if httpPort == "" {
		httpPort = "8080"
	}

	httpSrv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: httpServer,
	}

	go func() {
		logger.Info("HTTP Server listening on port "+httpPort, "server_startup", "port:"+httpPort)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panic(fmt.Sprintf("HTTP server error: %v", err), "http_server", "error")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("Shutdown signal received, shutting down gracefully...", "", "")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logger.Error(fmt.Sprintf("HTTP server shutdown error: %v", err), "", "")
	}

	logger.Info("Closing database connection...", "", "")
	db.Close()

	logger.Info("Servers shut down successfully", "", "")
}
