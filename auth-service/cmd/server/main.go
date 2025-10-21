package main

import (
	"context"

	"github.com/HatefBarari/microblog-auth/internal/infrastructure"
	"github.com/HatefBarari/microblog-auth/internal/presenter"
	"github.com/HatefBarari/microblog-auth/internal/repository"
	"github.com/HatefBarari/microblog-auth/internal/usecase"
	"github.com/HatefBarari/microblog-shared/pkg/email"
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
	log, err := logger.NewFile(cfg.Log.Level, "logs/auth.log")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	if err := mongo.Connect(cfg.Mongo.URI, cfg.Mongo.DBName, log); err != nil {
		panic(err)
	}

	// unique index on email
	idx := mongo.UsersColl().Indexes()
	_, _ = idx.CreateOne(context.Background(), mongoOptions.IndexModel{
		Keys: bson.M{"email": 1}, Options: (&options.IndexOptions{}).SetUnique(true),
	})
	emailSender := email.NewSender(email.Config{
		Host: cfg.Email.SMTP.Host,
		Port: cfg.Email.SMTP.Port,
		User: cfg.Email.SMTP.User,
		Pass: cfg.Email.SMTP.Pass,
	})

	// wiring
	repo := repository.NewMongoUserRepo()
	ucCfg := &usecase.Config{
		AccessSecret:   cfg.Auth.AccessSecret,
		RefreshSecret:  cfg.Auth.RefreshSecret,
		AccessTTLMin:   cfg.Auth.AccessTTLMin,
		RefreshTTLHour: cfg.Auth.RefreshTTLHour,
		EmailFrom:      cfg.Email.From,
	}
	uc := usecase.NewUserUseCase(repo, emailSender, ucCfg, log)
	
	// email usecase
	emailCfg := &usecase.EmailConfig{
		FromEmail:     cfg.Email.From,
		BaseURL:       cfg.Server.BaseURL,
		TokenSecret:   cfg.Auth.AccessSecret,
		TokenTTLHours: cfg.Auth.RefreshTTLHour,
	}
	emailUC := usecase.NewEmailUseCase(repo, emailSender, emailCfg, log)
	
	handler := presenter.NewHTTPHandler(uc, emailUC)

	if err := infrastructure.StartEcho(cfg, log, handler); err != nil {
		log.Fatal("start server", zap.Error(err))
	}
}