package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/HatefBarari/microblog-auth/internal/domain"
	"github.com/HatefBarari/microblog-auth/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock repositories
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) Send(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func TestEmailUseCase_SendVerificationEmail(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		userID         string
		email          string
		mockSetup      func(*MockEmailSender)
		expectedError  string
	}{
		{
			name:   "successful verification email",
			userID: "user123",
			email:  "user@example.com",
			mockSetup: func(mockSender *MockEmailSender) {
				mockSender.On("Send", "user@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:   "email send failure",
			userID: "user123",
			email:  "user@example.com",
			mockSetup: func(mockSender *MockEmailSender) {
				mockSender.On("Send", "user@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(errors.New("smtp error"))
			},
			expectedError: "failed to send verification email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockSender := new(MockEmailSender)
			tt.mockSetup(mockSender)

			uc := usecase.NewEmailUseCase(
				mockRepo,
				mockSender,
				&usecase.EmailConfig{
					FromEmail:     "noreply@microblog.com",
					BaseURL:       "http://localhost:8081",
					TokenSecret:   "secret",
					TokenTTLHours: 24,
				},
				logger,
			)

			err := uc.SendVerificationEmail(context.Background(), tt.userID, tt.email)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockSender.AssertExpectations(t)
		})
	}
}

func TestEmailUseCase_VerifyEmail(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		token          string
		mockSetup      func(*MockUserRepository)
		expectedError  string
	}{
		{
			name:  "successful email verification",
			token: "valid_token",
			mockSetup: func(mockRepo *MockUserRepository) {
				user := &domain.User{
					ID:       "user123",
					Email:    "user@example.com",
					Verified: false,
				}
				mockRepo.On("GetByID", mock.Anything, "user123").Return(user, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.ID == "user123" && u.Verified == true
				})).Return(nil)
			},
			expectedError: "",
		},
		{
			name:  "invalid token",
			token: "invalid_token",
			mockSetup: func(mockRepo *MockUserRepository) {
				// No repository calls expected for invalid token
			},
			expectedError: "invalid verification token",
		},
		{
			name:  "user not found",
			token: "valid_token",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByID", mock.Anything, "user123").Return(nil, nil)
			},
			expectedError: "user not found",
		},
		{
			name:  "email already verified",
			token: "valid_token",
			mockSetup: func(mockRepo *MockUserRepository) {
				user := &domain.User{
					ID:       "user123",
					Email:    "user@example.com",
					Verified: true, // already verified
				}
				mockRepo.On("GetByID", mock.Anything, "user123").Return(user, nil)
			},
			expectedError: "email already verified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockSender := new(MockEmailSender)
			tt.mockSetup(mockRepo)

			uc := usecase.NewEmailUseCase(
				mockRepo,
				mockSender,
				&usecase.EmailConfig{
					FromEmail:     "noreply@microblog.com",
					BaseURL:       "http://localhost:8081",
					TokenSecret:   "secret",
					TokenTTLHours: 24,
				},
				logger,
			)

			err := uc.VerifyEmail(context.Background(), tt.token)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEmailUseCase_SendPasswordResetEmail(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		email          string
		mockSetup      func(*MockUserRepository, *MockEmailSender)
		expectedError  string
	}{
		{
			name:  "successful password reset email",
			email: "user@example.com",
			mockSetup: func(mockRepo *MockUserRepository, mockSender *MockEmailSender) {
				user := &domain.User{
					ID:    "user123",
					Email: "user@example.com",
				}
				mockRepo.On("GetByEmail", mock.Anything, "user@example.com").Return(user, nil)
				mockSender.On("Send", "user@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:  "user not found (should not reveal)",
			email: "nonexistent@example.com",
			mockSetup: func(mockRepo *MockUserRepository, mockSender *MockEmailSender) {
				mockRepo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, nil)
				// No email should be sent
			},
			expectedError: "",
		},
		{
			name:  "email send failure",
			email: "user@example.com",
			mockSetup: func(mockRepo *MockUserRepository, mockSender *MockEmailSender) {
				user := &domain.User{
					ID:    "user123",
					Email: "user@example.com",
				}
				mockRepo.On("GetByEmail", mock.Anything, "user@example.com").Return(user, nil)
				mockSender.On("Send", "user@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
					Return(errors.New("smtp error"))
			},
			expectedError: "failed to send password reset email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockSender := new(MockEmailSender)
			tt.mockSetup(mockRepo, mockSender)

			uc := usecase.NewEmailUseCase(
				mockRepo,
				mockSender,
				&usecase.EmailConfig{
					FromEmail:     "noreply@microblog.com",
					BaseURL:       "http://localhost:8081",
					TokenSecret:   "secret",
					TokenTTLHours: 24,
				},
				logger,
			)

			err := uc.SendPasswordResetEmail(context.Background(), tt.email)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockSender.AssertExpectations(t)
		})
	}
}

func TestEmailUseCase_ResetPassword(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		token          string
		newPassword    string
		mockSetup      func(*MockUserRepository)
		expectedError  string
	}{
		{
			name:        "successful password reset",
			token:       "valid_token",
			newPassword: "newpassword123",
			mockSetup: func(mockRepo *MockUserRepository) {
				user := &domain.User{
					ID:           "user123",
					Email:        "user@example.com",
					PasswordHash: "old_hash",
				}
				mockRepo.On("GetByID", mock.Anything, "user123").Return(user, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.ID == "user123" && u.PasswordHash != "old_hash"
				})).Return(nil)
			},
			expectedError: "",
		},
		{
			name:        "invalid token",
			token:       "invalid_token",
			newPassword: "newpassword123",
			mockSetup: func(mockRepo *MockUserRepository) {
				// No repository calls expected for invalid token
			},
			expectedError: "invalid reset token",
		},
		{
			name:        "user not found",
			token:       "valid_token",
			newPassword: "newpassword123",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByID", mock.Anything, "user123").Return(nil, nil)
			},
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockSender := new(MockEmailSender)
			tt.mockSetup(mockRepo)

			uc := usecase.NewEmailUseCase(
				mockRepo,
				mockSender,
				&usecase.EmailConfig{
					FromEmail:     "noreply@microblog.com",
					BaseURL:       "http://localhost:8081",
					TokenSecret:   "secret",
					TokenTTLHours: 24,
				},
				logger,
			)

			err := uc.ResetPassword(context.Background(), tt.token, tt.newPassword)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
