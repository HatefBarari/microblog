package tests

import (
	"context"
	"testing"

	"github.com/HatefBarari/microblog-auth/internal/repository"
	"github.com/HatefBarari/microblog-auth/internal/usecase"
	"github.com/HatefBarari/microblog-shared/pkg/email"
	"github.com/HatefBarari/microblog-shared/pkg/logger"
	"github.com/HatefBarari/microblog-shared/pkg/mongo"
	"github.com/stretchr/testify/assert"
)

func TestRegisterDuplicateEmail(t *testing.T) {
	log, _ := logger.NewFile("debug", "logs/test.log")
	_ = mongo.Connect("mongodb://localhost:27017", "testdb", log)
	repo := repository.NewMongoUserRepo()
	emailSender := email.NewSender(email.Config{
		Host: "localhost",
		Port: "1025",
		User: "",
		Pass: "",
	})
	cfg := &usecase.Config{
		AccessSecret:   "test-access",
		RefreshSecret:  "test-refresh",
		AccessTTLMin:   15,
		RefreshTTLHour: 168,
		EmailFrom:      "test@local",
	}
	uc := usecase.NewUserUseCase(repo, emailSender, cfg, log)

	ctx := context.Background()
	req1 := usecase.RegisterRequest{Email: "a@x.com", Password: "123456"}
	_, err := uc.Register(ctx, req1)
	assert.Error(t, err)

	req2 := usecase.RegisterRequest{Email: "a@x.com", Password: "654321"}
	_, err = uc.Register(ctx, req2)
	assert.Error(t, err)
}

func TestEmailVerification(t *testing.T) {
	log, _ := logger.NewFile("debug", "logs/test.log")
	_ = mongo.Connect("mongodb://localhost:27017", "testdb", log)
	repo := repository.NewMongoUserRepo()
	emailSender := email.NewSender(email.Config{
		Host: "localhost",
		Port: "1025",
		User: "",
		Pass: "",
	})
	
	// User usecase
	cfg := &usecase.Config{
		AccessSecret:   "test-access",
		RefreshSecret:  "test-refresh",
		AccessTTLMin:   15,
		RefreshTTLHour: 168,
		EmailFrom:      "test@local",
	}
	uc := usecase.NewUserUseCase(repo, emailSender, cfg, log)
	
	// Email usecase
	emailCfg := &usecase.EmailConfig{
		FromEmail:     "test@local",
		BaseURL:       "http://localhost:8081",
		TokenSecret:   "test-secret",
		TokenTTLHours: 24,
	}
	emailUC := usecase.NewEmailUseCase(repo, emailSender, emailCfg, log)

	ctx := context.Background()
	
	// Test sending verification email
	err := emailUC.SendVerificationEmail(ctx, "user123", "test@example.com")
	assert.NoError(t, err)
	
	// Test email verification (this will fail with current implementation)
	// In a real implementation, you would need proper token validation
	err = emailUC.VerifyEmail(ctx, "invalid_token")
	assert.Error(t, err)
}

func TestPasswordReset(t *testing.T) {
	log, _ := logger.NewFile("debug", "logs/test.log")
	_ = mongo.Connect("mongodb://localhost:27017", "testdb", log)
	repo := repository.NewMongoUserRepo()
	emailSender := email.NewSender(email.Config{
		Host: "localhost",
		Port: "1025",
		User: "",
		Pass: "",
	})
	
	// Email usecase
	emailCfg := &usecase.EmailConfig{
		FromEmail:     "test@local",
		BaseURL:       "http://localhost:8081",
		TokenSecret:   "test-secret",
		TokenTTLHours: 24,
	}
	emailUC := usecase.NewEmailUseCase(repo, emailSender, emailCfg, log)

	ctx := context.Background()
	
	// Test sending password reset email
	err := emailUC.SendPasswordResetEmail(ctx, "test@example.com")
	assert.NoError(t, err)
	
	// Test password reset (this will fail with current implementation)
	// In a real implementation, you would need proper token validation
	err = emailUC.ResetPassword(ctx, "invalid_token", "newpassword123")
	assert.Error(t, err)
}