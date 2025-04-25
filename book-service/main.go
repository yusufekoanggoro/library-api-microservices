package main

import (
	"book-service/config"
	"book-service/config/key"
	deliveryG "book-service/internal/delivery/grpc"
	deliveryH "book-service/internal/delivery/http"
	"book-service/internal/delivery/http/routes"
	"book-service/internal/repository"
	"book-service/internal/usecase"
	"book-service/pkg/database"
	"book-service/pkg/logger"
	"book-service/pkg/middleware"
	"book-service/pkg/redis"
	sharedDomain "book-service/pkg/shared/domain"
	"book-service/pkg/token"
	"book-service/proto/book"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	logger := logger.NewLogger("book-service", logrus.DebugLevel, os.Stdout)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute) // timeout to shutdown server and init configuration
	defer func() {
		cancel()
		if r := recover(); r != nil {
			// notify if this service failed to start (to logstash or slack)
			fmt.Println("Failed to start book service:", r)
			fmt.Printf("Stack trace: \n%s\n", debug.Stack())
		}
	}()

	db := &database.GormDatabase{}
	if err := db.Connect(cfg); err != nil {
		logger.Panic(fmt.Sprintf("Database connection error: %v", err), "db-error", "connection")
	}

	if err := db.AutoMigrate(
		&sharedDomain.Author{},
		&sharedDomain.Book{},
		&sharedDomain.Borrowing{},
		&sharedDomain.Category{},
		&sharedDomain.Recommendation{},
		&sharedDomain.User{},
	); err != nil {
		logger.Panic(fmt.Sprintf("Failed to perform migration: %v", err), "migration", "error")
	}

	redisService := redis.NewRedisService(cfg)
	err := redisService.Ping(ctx)
	if err != nil {
		logger.Panic(fmt.Sprintf("Unable to connect to Redis: %v", err), "redis_connection", "error")
	}
	logger.Info("Redis is running and connected successfully", "redis_connection", "status:success")

	locker := redis.NewRedisLocker(redisService.GetClient())

	privateKeyPath := "config/key/private_key.pem"
	publicKeyPath := "config/key/public_key.pem"

	privateKey, publicKey, err := key.LoadRSAKeys(privateKeyPath, publicKeyPath)
	if err != nil {
		logger.Panic(fmt.Sprintf("Failed to load RSA keys: %v", err), "load rsa", "error")
	}

	jwt := token.NewJWT(publicKey, privateKey)

	repo := repository.NewRepository(db.GetDB())

	bookUsecase := usecase.NewBookUsecase(repo)
	stockUsecase := usecase.NewStockUsecase(repo)
	categoryUsecase := usecase.NewCategoryUsecase(repo)
	authorUsecase := usecase.NewAuthorUsecase(repo)
	userUsecase := usecase.NewUserUsecase(repo)
	borrowingUsecase := usecase.NewBorrowingUseCase(repo, locker)
	recommendationUsecase := usecase.NewRecommendationUsecase(repo)

	bookHandler := deliveryH.NewBookHandler(bookUsecase)
	stockHandler := deliveryH.NewStockHandler(stockUsecase)
	borrowingHandler := deliveryH.NewBorrowingHandler(borrowingUsecase)
	recommendationHandler := deliveryH.NewRecommendationHandler(recommendationUsecase)

	grpcBookHandler := deliveryG.NewBookHandler(userUsecase, authorUsecase, categoryUsecase)

	// Start gRPC Server
	grpcPort := cfg.GetGRPCPort()
	if grpcPort == "" {
		grpcPort = "50051"
	}

	grpcLis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.Panic(fmt.Sprintf("Failed to listen: %v", err), "grpc_listen", "error")
	}
	grpcSrv := grpc.NewServer()

	book.RegisterBookServiceServer(grpcSrv, grpcBookHandler)

	go func() {
		logger.Info("gRPC Server listening on port :"+grpcPort, "grpc_listen", "info")
		if err := grpcSrv.Serve(grpcLis); err != nil {
			logger.Panic(fmt.Sprintf("gRPC server error: %v", err), "grpc_server", "error")
		}
	}()

	// Setup routes
	r := gin.Default()

	middleware := middleware.NewMiddleware()

	httpPort := cfg.GetHTTPPort()
	if httpPort == "" {
		httpPort = "8080"
	}

	httpSrv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: routes.SetupRoutes(r, cfg, bookHandler, borrowingHandler, recommendationHandler, stockHandler, logger, jwt, middleware),
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

	grpcSrv.GracefulStop()

	logger.Info("Closing database connection...", "", "")
	db.Close()

	logger.Info("Servers shut down successfully", "", "")
}
