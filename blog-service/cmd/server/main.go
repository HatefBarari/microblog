package main

import (
	"context"

	"github.com/HatefBarari/microblog-blog/internal/infrastructure"
	"github.com/HatefBarari/microblog-blog/internal/presenter"
	"github.com/HatefBarari/microblog-blog/internal/repository"
	"github.com/HatefBarari/microblog-blog/internal/usecase"
	"github.com/HatefBarari/microblog-shared/pkg/logger"
	"github.com/HatefBarari/microblog-shared/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongoOptions "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	cfg, err := infrastructure.Load("configs/config.yaml")
	if err != nil {
		panic(err)
	}
	log, err := logger.NewFile(cfg.Log.Level, "logs/blog.log")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	if err := mongo.Connect(cfg.Mongo.URI, cfg.Mongo.DBName, log); err != nil {
		panic(err)
	}

	// indexes
	idx := mongo.DB().Collection("articles").Indexes()
	_, _ = idx.CreateOne(context.Background(), mongoOptions.IndexModel{
		Keys: bson.M{"slug": 1}, Options: (&options.IndexOptions{}).SetUnique(true),
	})
	_, _ = idx.CreateOne(context.Background(), mongoOptions.IndexModel{
		Keys: bson.M{"author_id": 1, "created_at": -1},
	})
	_, _ = mongo.DB().Collection("categories").Indexes().CreateOne(context.Background(), mongoOptions.IndexModel{
		Keys: bson.M{"slug": 1}, Options: (&options.IndexOptions{}).SetUnique(true),
	})

	// wiring – فقط مقادیر رو پاس بده نه struct
	articleRepo := repository.NewMongoArticleRepo()
	categoryRepo := repository.NewMongoCategoryRepo()
	commentRepo := repository.NewMongoCommentRepo()
	ratingRepo := repository.NewMongoRatingRepo()

	accessSecret := cfg.Auth.AccessToken
	refreshSecret := cfg.Auth.RefreshToken
	logLevel := cfg.Log.Level

	articleUC := usecase.NewArticleUseCase(articleRepo, ratingRepo, accessSecret, refreshSecret, logLevel, log)
	categoryUC := usecase.NewCategoryUseCase(categoryRepo, logLevel, log)
	commentUC := usecase.NewCommentUseCase(commentRepo, logLevel, log)
	ratingUC := usecase.NewRatingUseCase(ratingRepo, logLevel, log)

	handler := presenter.NewHTTPHandler(articleUC, categoryUC, commentUC, ratingUC)

	if err := infrastructure.StartEcho(cfg, log, handler); err != nil {
		log.Fatal("start server", zap.Error(err))
	}
}