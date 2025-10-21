package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/HatefBarari/microblog-auth/internal/domain"
	"github.com/HatefBarari/microblog-shared/pkg/auth"
	"github.com/HatefBarari/microblog-shared/pkg/email"
	"go.uber.org/zap"
)

type UserUseCase struct {
	repo   domain.UserRepository
	email  *email.Sender
	cfg    *Config
	log    *zap.Logger
}

type Config struct {
	AccessSecret   string
	RefreshSecret  string
	AccessTTLMin   int
	RefreshTTLHour int
	EmailFrom      string
}

func NewUserUseCase(repo domain.UserRepository, email *email.Sender, cfg *Config, log *zap.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, email: email, cfg: cfg, log: log}
}

func (uc *UserUseCase) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	// duplicate email?
	existing, _ := uc.repo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}
	// hash password
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	// create user
	u := &domain.User{
		Email:        req.Email,
		PasswordHash: hash,
		Role:         domain.RoleUser,
		Verified:     false,
		CreatedAt:    time.Now(),
	}
	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	// generate tokens
	acc, ref, err := auth.GenerateTokens(
		u.ID, string(u.Role),
		uc.cfg.AccessSecret, uc.cfg.RefreshSecret,
		uc.cfg.AccessTTLMin, uc.cfg.RefreshTTLHour,
	)
	if err != nil {
		return nil, err
	}
	// send verification email
	go func() {
		if err := uc.email.Send(u.Email, "Verify your account", "Click here: ..."); err != nil {
			uc.log.Error("send email", zap.Error(err))
		}
	}()
	return &RegisterResponse{
		AccessToken:  acc,
		RefreshToken: ref,
	}, nil
}

func (uc *UserUseCase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	u, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if u == nil || !auth.CheckPassword(u.PasswordHash, req.Password) {
		return nil, errors.New("invalid credentials")
	}
	if !u.Verified {
		return nil, errors.New("account not verified")
	}
	acc, ref, err := auth.GenerateTokens(
		u.ID, string(u.Role),
		uc.cfg.AccessSecret, uc.cfg.RefreshSecret,
		uc.cfg.AccessTTLMin, uc.cfg.RefreshTTLHour,
	)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{
		AccessToken:  acc,
		RefreshToken: ref,
	}, nil
}

func (uc *UserUseCase) Refresh(ctx context.Context, refreshToken string) (*RefreshResponse, error) {
	claims, err := auth.ValidateToken(refreshToken, uc.cfg.RefreshSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	// generate new access only
	acc, _, err := auth.GenerateTokens(
		claims.UserID, claims.Role,
		uc.cfg.AccessSecret, uc.cfg.RefreshSecret,
		uc.cfg.AccessTTLMin, uc.cfg.RefreshTTLHour,
	)
	if err != nil {
		return nil, err
	}
	return &RefreshResponse{AccessToken: acc}, nil
}