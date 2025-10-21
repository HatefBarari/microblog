package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HatefBarari/microblog-blog/internal/domain"
	"github.com/HatefBarari/microblog-blog/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock repositories
type MockArticleRepository struct {
	mock.Mock
}

func (m *MockArticleRepository) Create(ctx context.Context, article *domain.Article) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockArticleRepository) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Article), args.Error(1)
}

func (m *MockArticleRepository) GetBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Article), args.Error(1)
}

func (m *MockArticleRepository) List(ctx context.Context, filter domain.ListFilter) ([]*domain.Article, int, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*domain.Article), args.Int(1), args.Error(2)
}

func (m *MockArticleRepository) Update(ctx context.Context, article *domain.Article) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockArticleRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockArticleRepository) UpdateViewCount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

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

func TestArticleUseCase_Create(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		userID         string
		request        usecase.CreateArticleRequest
		mockSetup      func(*MockArticleRepository)
		expectedError  string
		expectedResult func(*usecase.ArticleResponse) bool
	}{
		{
			name:   "successful article creation",
			userID: "user123",
			request: usecase.CreateArticleRequest{
				Title:      "Test Article",
				Content:    "This is a test article content",
				Summary:    "Test summary",
				CategoryID: "cat123",
				Tags:       []string{"test", "go"},
				CoverURL:   "https://example.com/cover.jpg",
			},
			mockSetup: func(mockRepo *MockArticleRepository) {
				mockRepo.On("GetBySlug", mock.Anything, "test-article").Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(article *domain.Article) bool {
					return article.Title == "Test Article" &&
						article.AuthorID == "user123" &&
						article.Slug == "test-article" &&
						article.Status == domain.StatusDraft
				})).Return(nil)
			},
			expectedError: "",
			expectedResult: func(resp *usecase.ArticleResponse) bool {
				return resp.Title == "Test Article" &&
					resp.AuthorID == "user123" &&
					resp.Slug == "test-article" &&
					resp.Status == string(domain.StatusDraft)
			},
		},
		{
			name:   "too many tags",
			userID: "user123",
			request: usecase.CreateArticleRequest{
				Title:      "Test Article",
				Content:    "This is a test article content",
				Summary:    "Test summary",
				CategoryID: "cat123",
				Tags:       []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6", "tag7", "tag8", "tag9", "tag10", "tag11"},
			},
			mockSetup:     func(mockRepo *MockArticleRepository) {},
			expectedError: "max 10 tags",
		},
		{
			name:   "repository error",
			userID: "user123",
			request: usecase.CreateArticleRequest{
				Title:      "Test Article",
				Content:    "This is a test article content",
				Summary:    "Test summary",
				CategoryID: "cat123",
			},
			mockSetup: func(mockRepo *MockArticleRepository) {
				mockRepo.On("GetBySlug", mock.Anything, "test-article").Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockArticleRepository)
			mockRatingRepo := new(MockRatingRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewArticleUseCase(
				mockRepo,
				mockRatingRepo,
				"access_secret",
				"refresh_secret",
				"debug",
				logger,
			)

			result, err := uc.Create(context.Background(), tt.userID, tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.True(t, tt.expectedResult(result))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestArticleUseCase_GetBySlug(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		slug           string
		mockSetup      func(*MockArticleRepository)
		expectedError  string
		expectedResult func(*usecase.ArticleResponse) bool
	}{
		{
			name: "successful get by slug",
			slug: "test-article",
			mockSetup: func(mockRepo *MockArticleRepository) {
				article := &domain.Article{
					ID:        "article123",
					AuthorID:  "user123",
					Title:     "Test Article",
					Slug:      "test-article",
					Summary:   "Test summary",
					Content:   "Test content",
					CoverURL:  "https://example.com/cover.jpg",
					Status:    domain.StatusApproved,
					ViewCount: 10,
					RatingAvg: 4.5,
					CreatedAt: time.Now(),
				}
				mockRepo.On("GetBySlug", mock.Anything, "test-article").Return(article, nil)
				mockRepo.On("UpdateViewCount", mock.Anything, "article123").Return(nil)
			},
			expectedError: "",
			expectedResult: func(resp *usecase.ArticleResponse) bool {
				return resp.ID == "article123" &&
					resp.Title == "Test Article" &&
					resp.ViewCount == 11 // incremented
			},
		},
		{
			name: "article not found",
			slug: "non-existent",
			mockSetup: func(mockRepo *MockArticleRepository) {
				mockRepo.On("GetBySlug", mock.Anything, "non-existent").Return(nil, nil)
			},
			expectedError: "article not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockArticleRepository)
			mockRatingRepo := new(MockRatingRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewArticleUseCase(
				mockRepo,
				mockRatingRepo,
				"access_secret",
				"refresh_secret",
				"debug",
				logger,
			)

			result, err := uc.GetBySlug(context.Background(), tt.slug)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.True(t, tt.expectedResult(result))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestArticleUseCase_Update(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		userID         string
		articleID      string
		request        usecase.CreateArticleRequest
		mockSetup      func(*MockArticleRepository)
		expectedError  string
	}{
		{
			name:      "successful update",
			userID:    "user123",
			articleID: "article123",
			request: usecase.CreateArticleRequest{
				Title:      "Updated Article",
				Content:    "Updated content",
				Summary:    "Updated summary",
				CategoryID: "cat456",
				Tags:       []string{"updated", "test"},
			},
			mockSetup: func(mockRepo *MockArticleRepository) {
				existingArticle := &domain.Article{
					ID:       "article123",
					AuthorID: "user123",
					Title:    "Old Title",
					Slug:     "old-slug",
				}
				mockRepo.On("GetByID", mock.Anything, "article123").Return(existingArticle, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(article *domain.Article) bool {
					return article.Title == "Updated Article" &&
						article.AuthorID == "user123"
				})).Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "article not found",
			userID:    "user123",
			articleID: "nonexistent",
			request: usecase.CreateArticleRequest{
				Title:      "Updated Article",
				Content:    "Updated content",
				Summary:    "Updated summary",
				CategoryID: "cat456",
			},
			mockSetup: func(mockRepo *MockArticleRepository) {
				mockRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedError: "article not found",
		},
		{
			name:      "forbidden - not author",
			userID:    "user456",
			articleID: "article123",
			request: usecase.CreateArticleRequest{
				Title:      "Updated Article",
				Content:    "Updated content",
				Summary:    "Updated summary",
				CategoryID: "cat456",
			},
			mockSetup: func(mockRepo *MockArticleRepository) {
				existingArticle := &domain.Article{
					ID:       "article123",
					AuthorID: "user123", // different author
					Title:    "Old Title",
				}
				mockRepo.On("GetByID", mock.Anything, "article123").Return(existingArticle, nil)
			},
			expectedError: "forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockArticleRepository)
			mockRatingRepo := new(MockRatingRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewArticleUseCase(
				mockRepo,
				mockRatingRepo,
				"access_secret",
				"refresh_secret",
				"debug",
				logger,
			)

			result, err := uc.Update(context.Background(), tt.userID, tt.articleID, tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "Updated Article", result.Title)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestArticleUseCase_Delete(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name          string
		userID        string
		articleID     string
		mockSetup     func(*MockArticleRepository)
		expectedError string
	}{
		{
			name:      "successful delete",
			userID:    "user123",
			articleID: "article123",
			mockSetup: func(mockRepo *MockArticleRepository) {
				existingArticle := &domain.Article{
					ID:       "article123",
					AuthorID: "user123",
				}
				mockRepo.On("GetByID", mock.Anything, "article123").Return(existingArticle, nil)
				mockRepo.On("Delete", mock.Anything, "article123").Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "article not found",
			userID:    "user123",
			articleID: "nonexistent",
			mockSetup: func(mockRepo *MockArticleRepository) {
				mockRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedError: "article not found",
		},
		{
			name:      "forbidden - not author",
			userID:    "user456",
			articleID: "article123",
			mockSetup: func(mockRepo *MockArticleRepository) {
				existingArticle := &domain.Article{
					ID:       "article123",
					AuthorID: "user123", // different author
				}
				mockRepo.On("GetByID", mock.Anything, "article123").Return(existingArticle, nil)
			},
			expectedError: "forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockArticleRepository)
			mockRatingRepo := new(MockRatingRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewArticleUseCase(
				mockRepo,
				mockRatingRepo,
				"access_secret",
				"refresh_secret",
				"debug",
				logger,
			)

			err := uc.Delete(context.Background(), tt.userID, tt.articleID)

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
