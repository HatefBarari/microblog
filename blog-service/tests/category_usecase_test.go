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

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) List(ctx context.Context) ([]*domain.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) ListByParent(ctx context.Context, parentID string) ([]*domain.Category, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func TestCategoryUseCase_Create(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		categoryName   string
		parentID       string
		mockSetup      func(*MockCategoryRepository)
		expectedError  string
		expectedResult func(*usecase.CategoryResponse) bool
	}{
		{
			name:         "successful category creation",
			categoryName: "Technology",
			parentID:     "",
			mockSetup: func(mockRepo *MockCategoryRepository) {
				mockRepo.On("GetBySlug", mock.Anything, "technology").Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(category *domain.Category) bool {
					return category.Name == "Technology" &&
						category.Slug == "technology" &&
						category.ParentID == ""
				})).Return(nil)
			},
			expectedError: "",
			expectedResult: func(resp *usecase.CategoryResponse) bool {
				return resp.Name == "Technology" &&
					resp.Slug == "technology" &&
					resp.ParentID == ""
			},
		},
		{
			name:         "successful subcategory creation",
			categoryName: "Web Development",
			parentID:     "parent123",
			mockSetup: func(mockRepo *MockCategoryRepository) {
				mockRepo.On("GetBySlug", mock.Anything, "web-development").Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(category *domain.Category) bool {
					return category.Name == "Web Development" &&
						category.Slug == "web-development" &&
						category.ParentID == "parent123"
				})).Return(nil)
			},
			expectedError: "",
			expectedResult: func(resp *usecase.CategoryResponse) bool {
				return resp.Name == "Web Development" &&
					resp.Slug == "web-development" &&
					resp.ParentID == "parent123"
			},
		},
		{
			name:         "duplicate slug",
			categoryName: "Technology",
			parentID:     "",
			mockSetup: func(mockRepo *MockCategoryRepository) {
				existingCategory := &domain.Category{
					ID:   "existing123",
					Name: "Technology",
					Slug: "technology",
				}
				mockRepo.On("GetBySlug", mock.Anything, "technology").Return(existingCategory, nil)
			},
			expectedError: "category with this slug already exists",
		},
		{
			name:         "repository error",
			categoryName: "Technology",
			parentID:     "",
			mockSetup: func(mockRepo *MockCategoryRepository) {
				mockRepo.On("GetBySlug", mock.Anything, "technology").Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCategoryRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewCategoryUseCase(mockRepo, logger)

			result, err := uc.Create(context.Background(), tt.categoryName, tt.parentID)

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

func TestCategoryUseCase_ListTree(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		mockSetup      func(*MockCategoryRepository)
		expectedError  string
		expectedResult func([]*usecase.CategoryResponse) bool
	}{
		{
			name: "successful tree listing",
			mockSetup: func(mockRepo *MockCategoryRepository) {
				categories := []*domain.Category{
					{ID: "cat1", Name: "Technology", Slug: "technology", ParentID: ""},
					{ID: "cat2", Name: "Programming", Slug: "programming", ParentID: "cat1"},
					{ID: "cat3", Name: "Science", Slug: "science", ParentID: ""},
				}
				mockRepo.On("List", mock.Anything).Return(categories, nil)
			},
			expectedError: "",
			expectedResult: func(resp []*usecase.CategoryResponse) bool {
				return len(resp) == 2 && // 2 root categories
					resp[0].Name == "Technology" &&
					len(resp[0].Children) == 1 && // 1 child
					resp[0].Children[0].Name == "Programming" &&
					resp[1].Name == "Science" &&
					len(resp[1].Children) == 0 // no children
			},
		},
		{
			name: "empty tree",
			mockSetup: func(mockRepo *MockCategoryRepository) {
				mockRepo.On("List", mock.Anything).Return([]*domain.Category{}, nil)
			},
			expectedError: "",
			expectedResult: func(resp []*usecase.CategoryResponse) bool {
				return len(resp) == 0
			},
		},
		{
			name: "repository error",
			mockSetup: func(mockRepo *MockCategoryRepository) {
				mockRepo.On("List", mock.Anything).Return([]*domain.Category{}, errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCategoryRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewCategoryUseCase(mockRepo, logger)

			result, err := uc.ListTree(context.Background())

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
