package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HatefBarari/microblog-media/internal/infrastructure"
	"github.com/HatefBarari/microblog-media/internal/presenter"
	"github.com/HatefBarari/microblog-media/internal/repository"
	"github.com/HatefBarari/microblog-media/internal/usecase"
	"github.com/HatefBarari/microblog-shared/pkg/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	config, err := infrastructure.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	logger, err := infrastructure.NewLogger(&config.Log)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting media service", 
		zap.String("port", config.Server.Port),
		zap.String("host", config.Server.Host))

	// Connect to MongoDB
	mongoClient, err := mongo.NewClient(config.Database.URI)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(config.Database.Database)

	// Initialize repositories
	mediaRepo := repository.NewMongoMediaRepository(db)

	// Initialize storage
	storage := infrastructure.NewFileStorage(
		config.Storage.BasePath,
		config.Storage.BaseURL,
		logger,
	)

	// Initialize use cases
	mediaUC := usecase.NewMediaUseCase(
		mediaRepo,
		storage,
		&usecase.Config{
			MaxFileSize:   config.Media.MaxFileSize,
			AllowedTypes:  config.Media.AllowedTypes,
			StoragePath:   config.Storage.BasePath,
			BaseURL:       config.Storage.BaseURL,
			ThumbnailSize: config.Media.ThumbnailSize,
		},
		logger,
	)

	// Initialize handlers
	handler := presenter.NewHTTPHandler(mediaUC)

	// Initialize server
	server := infrastructure.NewEchoServer(config)
	server.SetupRoutes(handler)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Media service started successfully")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down media service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(); err != nil {
		logger.Error("Failed to shutdown server", zap.Error(err))
	}

	logger.Info("Media service stopped")
}
