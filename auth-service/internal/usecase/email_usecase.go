package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HatefBarari/microblog-auth/internal/domain"
	"github.com/HatefBarari/microblog-shared/pkg/email"
	"go.uber.org/zap"
)

type EmailUseCase struct {
	repo   domain.UserRepository
	email  *email.Sender
	cfg    *EmailConfig
	log    *zap.Logger
}

type EmailConfig struct {
	FromEmail     string
	BaseURL       string
	TokenSecret   string
	TokenTTLHours int
}

func NewEmailUseCase(repo domain.UserRepository, emailSender *email.Sender, cfg *EmailConfig, log *zap.Logger) *EmailUseCase {
	return &EmailUseCase{
		repo:  repo,
		email: emailSender,
		cfg:   cfg,
		log:   log,
	}
}

// SendVerificationEmail sends verification email to user
func (uc *EmailUseCase) SendVerificationEmail(ctx context.Context, userID, email string) error {
	// Generate verification token
	token, err := uc.generateVerificationToken(userID)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Create verification URL
	verificationURL := fmt.Sprintf("%s/verify?token=%s", uc.cfg.BaseURL, token)

	// Create email content
	subject := "تایید حساب کاربری - Microblog"
	body := fmt.Sprintf(`
سلام،

برای تایید حساب کاربری خود روی لینک زیر کلیک کنید:

%s

این لینک تا %d ساعت معتبر است.

با تشکر،
تیم Microblog
`, verificationURL, uc.cfg.TokenTTLHours)

	// Send email
	if err := uc.email.Send(email, subject, body); err != nil {
		uc.log.Error("failed to send verification email", 
			zap.String("user_id", userID),
			zap.String("email", email),
			zap.Error(err))
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	uc.log.Info("verification email sent successfully",
		zap.String("user_id", userID),
		zap.String("email", email))

	return nil
}

// VerifyEmail verifies user email with token
func (uc *EmailUseCase) VerifyEmail(ctx context.Context, token string) error {
	// Validate token
	userID, err := uc.validateVerificationToken(token)
	if err != nil {
		return fmt.Errorf("invalid verification token: %w", err)
	}

	// Get user
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Check if already verified
	if user.Verified {
		return errors.New("email already verified")
	}

	// Update user verification status
	user.Verified = true
	if err := uc.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user verification: %w", err)
	}

	uc.log.Info("email verified successfully",
		zap.String("user_id", userID),
		zap.String("email", user.Email))

	return nil
}

// SendPasswordResetEmail sends password reset email
func (uc *EmailUseCase) SendPasswordResetEmail(ctx context.Context, email string) error {
	// Get user by email
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		// Don't reveal if email exists or not
		uc.log.Info("password reset requested for non-existent email",
			zap.String("email", email))
		return nil
	}

	// Generate reset token
	token, err := uc.generatePasswordResetToken(user.ID)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Create reset URL
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", uc.cfg.BaseURL, token)

	// Create email content
	subject := "بازیابی رمز عبور - Microblog"
	body := fmt.Sprintf(`
سلام،

برای بازیابی رمز عبور خود روی لینک زیر کلیک کنید:

%s

این لینک تا %d ساعت معتبر است.

اگر شما این درخواست را نکرده‌اید، این ایمیل را نادیده بگیرید.

با تشکر،
تیم Microblog
`, resetURL, uc.cfg.TokenTTLHours)

	// Send email
	if err := uc.email.Send(email, subject, body); err != nil {
		uc.log.Error("failed to send password reset email",
			zap.String("user_id", user.ID),
			zap.String("email", email),
			zap.Error(err))
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	uc.log.Info("password reset email sent successfully",
		zap.String("user_id", user.ID),
		zap.String("email", email))

	return nil
}

// ResetPassword resets user password with token
func (uc *EmailUseCase) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate token
	userID, err := uc.validatePasswordResetToken(token)
	if err != nil {
		return fmt.Errorf("invalid reset token: %w", err)
	}

	// Get user
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Hash new password
	hash, err := uc.hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	user.PasswordHash = hash
	if err := uc.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	uc.log.Info("password reset successfully",
		zap.String("user_id", userID))

	return nil
}

// generateVerificationToken generates a verification token for user
func (uc *EmailUseCase) generateVerificationToken(userID string) (string, error) {
	// In a real implementation, you would use JWT or a secure token
	// For now, we'll use a simple approach
	token := fmt.Sprintf("verify_%s_%d", userID, time.Now().Unix())
	return token, nil
}

// validateVerificationToken validates a verification token
func (uc *EmailUseCase) validateVerificationToken(token string) (string, error) {
	// In a real implementation, you would validate JWT or check database
	// For now, we'll use a simple approach
	if len(token) < 10 {
		return "", errors.New("invalid token format")
	}
	
	// Extract user ID from token (simplified)
	// In production, use proper JWT validation
	userID := "user123" // This should be extracted from token
	return userID, nil
}

// generatePasswordResetToken generates a password reset token
func (uc *EmailUseCase) generatePasswordResetToken(userID string) (string, error) {
	// In a real implementation, you would use JWT or a secure token
	token := fmt.Sprintf("reset_%s_%d", userID, time.Now().Unix())
	return token, nil
}

// validatePasswordResetToken validates a password reset token
func (uc *EmailUseCase) validatePasswordResetToken(token string) (string, error) {
	// In a real implementation, you would validate JWT or check database
	if len(token) < 10 {
		return "", errors.New("invalid token format")
	}
	
	// Extract user ID from token (simplified)
	userID := "user123" // This should be extracted from token
	return userID, nil
}

// hashPassword hashes a password
func (uc *EmailUseCase) hashPassword(password string) (string, error) {
	// Import the hash function from shared package
	// This is a placeholder - you should use the actual hash function
	return "hashed_" + password, nil
}