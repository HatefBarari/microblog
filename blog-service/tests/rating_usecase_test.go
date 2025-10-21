package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/HatefBarari/microblog-blog/internal/domain"
	"github.com/HatefBarari/microblog-blog/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockRatingRepository struct {
	mock.Mock
}

func (m *MockRatingRepository) Create(ctx context.Context, rating *domain.Rating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}

func (m *MockRatingRepository) GetByUserAndTarget(ctx context.Context, userID, targetID, targetType string) (*domain.Rating, error) {
	args := m.Called(ctx, userID, targetID, targetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Rating), args.Error(1)
}

func (m *MockRatingRepository) Update(ctx context.Context, rating *domain.Rating) error {
	args := m.Called(ctx, rating)
	return args.Error(0)
}

func (m *MockRatingRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRatingRepository) GetAverageByTarget(ctx context.Context, targetID, targetType string) (float64, error) {
	args := m.Called(ctx, targetID, targetType)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRatingRepository) UpdateTargetAverage(ctx context.Context, targetID, targetType string, avg float64) error {
	args := m.Called(ctx, targetID, targetType, avg)
	return args.Error(0)
}

func TestRatingUseCase_RateArticle(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		userID         string
		articleID      string
		stars          int
		mockSetup      func(*MockRatingRepository)
		expectedError  string
	}{
		{
			name:      "successful new rating",
			userID:    "user123",
			articleID: "article123",
			stars:     5,
			mockSetup: func(mockRepo *MockRatingRepository) {
				mockRepo.On("GetByUserAndTarget", mock.Anything, "user123", "article123", "article").
					Return(nil, nil) // no existing rating
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(rating *domain.Rating) bool {
					return rating.UserID == "user123" &&
						rating.TargetID == "article123" &&
						rating.Type == "article" &&
						rating.Stars == 5
				})).Return(nil)
				mockRepo.On("GetAverageByTarget", mock.Anything, "article123", "article").
					Return(5.0, nil)
				mockRepo.On("UpdateTargetAverage", mock.Anything, "article123", "article", 5.0).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "successful rating update",
			userID:    "user123",
			articleID: "article123",
			stars:     4,
			mockSetup: func(mockRepo *MockRatingRepository) {
				existingRating := &domain.Rating{
					ID:       "rating123",
					UserID:   "user123",
					TargetID: "article123",
					Type:     "article",
					Stars:    3,
				}
				mockRepo.On("GetByUserAndTarget", mock.Anything, "user123", "article123", "article").
					Return(existingRating, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(rating *domain.Rating) bool {
					return rating.ID == "rating123" &&
						rating.Stars == 4
				})).Return(nil)
				mockRepo.On("GetAverageByTarget", mock.Anything, "article123", "article").
					Return(4.0, nil)
				mockRepo.On("UpdateTargetAverage", mock.Anything, "article123", "article", 4.0).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "invalid stars - too low",
			userID:    "user123",
			articleID: "article123",
			stars:     0,
			mockSetup: func(mockRepo *MockRatingRepository) {
				// No repository calls expected for invalid input
			},
			expectedError: "stars must be between 1 and 5",
		},
		{
			name:      "invalid stars - too high",
			userID:    "user123",
			articleID: "article123",
			stars:     6,
			mockSetup: func(mockRepo *MockRatingRepository) {
				// No repository calls expected for invalid input
			},
			expectedError: "stars must be between 1 and 5",
		},
		{
			name:      "repository error on create",
			userID:    "user123",
			articleID: "article123",
			stars:     5,
			mockSetup: func(mockRepo *MockRatingRepository) {
				mockRepo.On("GetByUserAndTarget", mock.Anything, "user123", "article123", "article").
					Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRatingRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewRatingUseCase(mockRepo, logger)

			err := uc.RateArticle(context.Background(), tt.userID, tt.articleID, tt.stars)

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

func TestRatingUseCase_Delete(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name          string
		userID        string
		targetID      string
		targetType    string
		mockSetup     func(*MockRatingRepository)
		expectedError string
	}{
		{
			name:       "successful rating deletion",
			userID:     "user123",
			targetID:   "article123",
			targetType: "article",
			mockSetup: func(mockRepo *MockRatingRepository) {
				existingRating := &domain.Rating{
					ID:       "rating123",
					UserID:   "user123",
					TargetID: "article123",
					Type:     "article",
					Stars:    5,
				}
				mockRepo.On("GetByUserAndTarget", mock.Anything, "user123", "article123", "article").
					Return(existingRating, nil)
				mockRepo.On("Delete", mock.Anything, "rating123").Return(nil)
				mockRepo.On("GetAverageByTarget", mock.Anything, "article123", "article").
					Return(0.0, nil)
				mockRepo.On("UpdateTargetAverage", mock.Anything, "article123", "article", 0.0).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:       "rating not found",
			userID:     "user123",
			targetID:   "article123",
			targetType: "article",
			mockSetup: func(mockRepo *MockRatingRepository) {
				mockRepo.On("GetByUserAndTarget", mock.Anything, "user123", "article123", "article").
					Return(nil, nil)
			},
			expectedError: "rating not found",
		},
		{
			name:       "repository error on get",
			userID:     "user123",
			targetID:   "article123",
			targetType: "article",
			mockSetup: func(mockRepo *MockRatingRepository) {
				mockRepo.On("GetByUserAndTarget", mock.Anything, "user123", "article123", "article").
					Return(nil, errors.New("database error"))
			},
			expectedError: "database error",
		},
		{
			name:       "repository error on delete",
			userID:     "user123",
			targetID:   "article123",
			targetType: "article",
			mockSetup: func(mockRepo *MockRatingRepository) {
				existingRating := &domain.Rating{
					ID:       "rating123",
					UserID:   "user123",
					TargetID: "article123",
					Type:     "article",
					Stars:    5,
				}
				mockRepo.On("GetByUserAndTarget", mock.Anything, "user123", "article123", "article").
					Return(existingRating, nil)
				mockRepo.On("Delete", mock.Anything, "rating123").Return(errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRatingRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewRatingUseCase(mockRepo, logger)

			err := uc.Delete(context.Background(), tt.userID, tt.targetID, tt.targetType)

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
