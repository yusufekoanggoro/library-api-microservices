package main

import (
	"auth-service/config"
	"auth-service/config/key"
	deliveryH "auth-service/internal/delivery/http"
	"auth-service/internal/delivery/http/routes"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"auth-service/pkg/database"
	"auth-service/pkg/logger"
	"auth-service/pkg/middleware"
	sharedDomain "auth-service/pkg/shared/domain"
	"auth-service/pkg/token"
	"auth-service/seeders"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"auth-service/internal/grpcservice"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute) // timeout to shutdown server and init configuration
	defer func() {
		cancel()
		if r := recover(); r != nil {
			// notify if this service failed to start (to logstash or slack)
			fmt.Println("Failed to start auth service:", r)
			fmt.Printf("Stack trace: \n%s\n", debug.Stack())
		}
	}()

	cfg := config.LoadConfig()

	logger := logger.NewLogger("auth-service", logrus.DebugLevel, os.Stdout)

	db := &database.GormDatabase{}
	if err := db.Connect(cfg); err != nil {
		logger.Panic(fmt.Sprintf("Database connection error: %v", err), "db-error", "connection")
	}

	if err := db.AutoMigrate(
		&sharedDomain.User{},
	); err != nil {
		logger.Error(fmt.Sprintf("Failed to perform migration: %v", err), "migration", "error")
	}

	var users []*sharedDomain.User
	if u, err := seeders.SeedUsers(db.GetDB()); err != nil {
		logger.Info(fmt.Sprintf("Error seeding users: %v", err), "seeders", "info")
	} else {
		users = u
		logger.Info("Seeding users completed successfully", "seeders", "success")
	}

	bookClient := grpcservice.NewBookGRPCClient(cfg.GetBookGRPCHost() + ":" + cfg.GetBookGRPCPort())
	for _, user := range users {
		_, err := bookClient.SaveUser(ctx, user.ToProto())
		if err != nil {
			logger.Error(err.Error(), "bookClientSaveUser.", "error")
		}
	}

	privateKeyPath := "config/key/private_key.pem"
	publicKeyPath := "config/key/public_key.pem"

	privateKey, publicKey, err := key.LoadRSAKeys(privateKeyPath, publicKeyPath)
	if err != nil {
		logger.Panic(fmt.Sprintf("Failed to load RSA keys: %v", err), "load rsa", "error")
	}
	token := token.NewJWT(publicKey, privateKey)

	// Setup repository, usecase, handler, routes
	userRepo := repository.NewUserRepository(db.GetDB())
	userUsecase := usecase.NewAuthorUsecase(userRepo, bookClient, logger)
	authUsecase := usecase.NewAuthUsecase(userRepo, token, bookClient, logger)

	userHandler := deliveryH.NeUserHandler(userUsecase)
	authHandler := deliveryH.NewAuthHandler(authUsecase)

	middleware := middleware.NewMiddleware()

	r := gin.Default()

	httpPort := cfg.GetHTTPPort()
	if httpPort == "" {
		httpPort = "8080"
	}

	httpSrv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: routes.SetupRoutes(r, cfg, authHandler, userHandler, logger, token, middleware),
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
