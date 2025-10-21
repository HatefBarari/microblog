package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HatefBarari/microblog-blog/internal/domain"
	"github.com/HatefBarari/microblog-blog/internal/presenter"
	"github.com/HatefBarari/microblog-blog/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock use cases
type MockArticleUseCase struct {
	mock.Mock
}

func (m *MockArticleUseCase) Create(ctx context.Context, userID string, req usecase.CreateArticleRequest) (*usecase.ArticleResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ArticleResponse), args.Error(1)
}

func (m *MockArticleUseCase) GetBySlug(ctx context.Context, slug string) (*usecase.ArticleResponse, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ArticleResponse), args.Error(1)
}

func (m *MockArticleUseCase) List(ctx context.Context, filter domain.ListFilter) ([]*usecase.ArticleResponse, int, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*usecase.ArticleResponse), args.Int(1), args.Error(2)
}

func (m *MockArticleUseCase) Update(ctx context.Context, userID string, id string, req usecase.CreateArticleRequest) (*usecase.ArticleResponse, error) {
	args := m.Called(ctx, userID, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ArticleResponse), args.Error(1)
}

func (m *MockArticleUseCase) Delete(ctx context.Context, userID string, id string) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}

type MockCategoryUseCase struct {
	mock.Mock
}

func (m *MockCategoryUseCase) ListTree(ctx context.Context) ([]*usecase.CategoryResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*usecase.CategoryResponse), args.Error(1)
}

func (m *MockCategoryUseCase) Create(ctx context.Context, name, parentID string) (*usecase.CategoryResponse, error) {
	args := m.Called(ctx, name, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.CategoryResponse), args.Error(1)
}

type MockCommentUseCase struct {
	mock.Mock
}

func (m *MockCommentUseCase) ListByArticle(ctx context.Context, articleID string, status domain.CommentStatus) ([]*usecase.CommentResponse, error) {
	args := m.Called(ctx, articleID, status)
	return args.Get(0).([]*usecase.CommentResponse), args.Error(1)
}

func (m *MockCommentUseCase) Create(ctx context.Context, userID, articleID string, req usecase.CreateCommentRequest) (*usecase.CommentResponse, error) {
	args := m.Called(ctx, userID, articleID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.CommentResponse), args.Error(1)
}

func (m *MockCommentUseCase) UpdateStatus(ctx context.Context, id string, status domain.CommentStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

type MockRatingUseCase struct {
	mock.Mock
}

func (m *MockRatingUseCase) RateArticle(ctx context.Context, userID, articleID string, stars int) error {
	args := m.Called(ctx, userID, articleID, stars)
	return args.Error(0)
}

func (m *MockRatingUseCase) Delete(ctx context.Context, userID, targetID, targetType string) error {
	args := m.Called(ctx, userID, targetID, targetType)
	return args.Error(0)
}

func setupTestHandler() (*presenter.HTTPHandler, *MockArticleUseCase, *MockCategoryUseCase, *MockCommentUseCase, *MockRatingUseCase) {
	mockArticleUC := new(MockArticleUseCase)
	mockCategoryUC := new(MockCategoryUseCase)
	mockCommentUC := new(MockCommentUseCase)
	mockRatingUC := new(MockRatingUseCase)
	
	handler := presenter.NewHTTPHandler(mockArticleUC, mockCategoryUC, mockCommentUC, mockRatingUC)
	
	return handler, mockArticleUC, mockCategoryUC, mockCommentUC, mockRatingUC
}

func TestHTTPHandler_CreateArticle(t *testing.T) {
	handler, mockArticleUC, _, _, _ := setupTestHandler()
	
	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		mockSetup      func(*MockArticleUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful article creation",
			userID: "user123",
			requestBody: usecase.CreateArticleRequest{
				Title:      "Test Article",
				Content:    "This is a test article content",
				Summary:    "Test summary",
				CategoryID: "cat123",
				Tags:       []string{"test", "go"},
				CoverURL:   "https://example.com/cover.jpg",
			},
			mockSetup: func(mockUC *MockArticleUseCase) {
				mockUC.On("Create", mock.Anything, "user123", mock.AnythingOfType("usecase.CreateArticleRequest")).
					Return(&usecase.ArticleResponse{
						ID:        "article123",
						AuthorID:  "user123",
						Title:     "Test Article",
						Slug:      "test-article",
						Status:    "draft",
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "invalid request body",
			userID: "user123",
			requestBody: map[string]interface{}{
				"title": "", // invalid - too short
			},
			mockSetup:      func(mockUC *MockArticleUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation error",
		},
		{
			name:   "use case error",
			userID: "user123",
			requestBody: usecase.CreateArticleRequest{
				Title:      "Test Article",
				Content:    "This is a test article content",
				Summary:    "Test summary",
				CategoryID: "cat123",
			},
			mockSetup: func(mockUC *MockArticleUseCase) {
				mockUC.On("Create", mock.Anything, "user123", mock.AnythingOfType("usecase.CreateArticleRequest")).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockArticleUC.ExpectedCalls = nil
			tt.mockSetup(mockArticleUC)

			// Setup request
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)

			// Execute
			err := handler.CreateArticle(c)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockArticleUC.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_GetArticleBySlug(t *testing.T) {
	handler, mockArticleUC, _, _, _ := setupTestHandler()
	
	tests := []struct {
		name           string
		slug           string
		mockSetup      func(*MockArticleUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful get by slug",
			slug: "test-article",
			mockSetup: func(mockUC *MockArticleUseCase) {
				mockUC.On("GetBySlug", mock.Anything, "test-article").
					Return(&usecase.ArticleResponse{
						ID:        "article123",
						Title:     "Test Article",
						Slug:      "test-article",
						ViewCount: 10,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "article not found",
			slug: "non-existent",
			mockSetup: func(mockUC *MockArticleUseCase) {
				mockUC.On("GetBySlug", mock.Anything, "non-existent").
					Return(nil, errors.New("article not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "article not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockArticleUC.ExpectedCalls = nil
			tt.mockSetup(mockArticleUC)

			// Setup request
			req := httptest.NewRequest(http.MethodGet, "/articles/"+tt.slug, nil)
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames("slug")
			c.SetParamValues(tt.slug)

			// Execute
			err := handler.GetArticleBySlug(c)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockArticleUC.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_ListArticles(t *testing.T) {
	handler, mockArticleUC, _, _, _ := setupTestHandler()
	
	tests := []struct {
		name           string
		mockSetup      func(*MockArticleUseCase)
		expectedStatus int
	}{
		{
			name: "successful list articles",
			mockSetup: func(mockUC *MockArticleUseCase) {
				articles := []*usecase.ArticleResponse{
					{ID: "article1", Title: "Article 1"},
					{ID: "article2", Title: "Article 2"},
				}
				mockUC.On("List", mock.Anything, mock.AnythingOfType("domain.ListFilter")).
					Return(articles, 2, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "use case error",
			mockSetup: func(mockUC *MockArticleUseCase) {
				mockUC.On("List", mock.Anything, mock.AnythingOfType("domain.ListFilter")).
					Return([]*usecase.ArticleResponse{}, 0, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockArticleUC.ExpectedCalls = nil
			tt.mockSetup(mockArticleUC)

			// Setup request
			req := httptest.NewRequest(http.MethodGet, "/articles", nil)
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)

			// Execute
			err := handler.ListArticles(c)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockArticleUC.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_UpdateArticle(t *testing.T) {
	handler, mockArticleUC, _, _, _ := setupTestHandler()
	
	tests := []struct {
		name           string
		userID         string
		articleID      string
		requestBody    interface{}
		mockSetup      func(*MockArticleUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "successful update",
			userID:    "user123",
			articleID: "article123",
			requestBody: usecase.CreateArticleRequest{
				Title:      "Updated Article",
				Content:    "Updated content",
				Summary:    "Updated summary",
				CategoryID: "cat456",
			},
			mockSetup: func(mockUC *MockArticleUseCase) {
				mockUC.On("Update", mock.Anything, "user123", "article123", mock.AnythingOfType("usecase.CreateArticleRequest")).
					Return(&usecase.ArticleResponse{
						ID:     "article123",
						Title:  "Updated Article",
						Status: "draft",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "article not found",
			userID:    "user123",
			articleID: "nonexistent",
			requestBody: usecase.CreateArticleRequest{
				Title:      "Updated Article",
				Content:    "Updated content",
				Summary:    "Updated summary",
				CategoryID: "cat456",
			},
			mockSetup: func(mockUC *MockArticleUseCase) {
				mockUC.On("Update", mock.Anything, "user123", "nonexistent", mock.AnythingOfType("usecase.CreateArticleRequest")).
					Return(nil, errors.New("article not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "article not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockArticleUC.ExpectedCalls = nil
			tt.mockSetup(mockArticleUC)

			// Setup request
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/articles/"+tt.articleID, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)
			c.SetParamNames("id")
			c.SetParamValues(tt.articleID)

			// Execute
			err := handler.UpdateArticle(c)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockArticleUC.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_DeleteArticle(t *testing.T) {
	handler, mockArticleUC, _, _, _ := setupTestHandler()
	
	tests := []struct {
		name           string
		userID         string
		articleID      string
		mockSetup      func(*MockArticleUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "successful delete",
			userID:    "user123",
			articleID: "article123",
			mockSetup: func(mockUC *MockArticleUseCase) {
				mockUC.On("Delete", mock.Anything, "user123", "article123").Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "article not found",
			userID:    "user123",
			articleID: "nonexistent",
			mockSetup: func(mockUC *MockArticleUseCase) {
				mockUC.On("Delete", mock.Anything, "user123", "nonexistent").
					Return(errors.New("article not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "article not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockArticleUC.ExpectedCalls = nil
			tt.mockSetup(mockArticleUC)

			// Setup request
			req := httptest.NewRequest(http.MethodDelete, "/articles/"+tt.articleID, nil)
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)
			c.SetParamNames("id")
			c.SetParamValues(tt.articleID)

			// Execute
			err := handler.DeleteArticle(c)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockArticleUC.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_CreateComment(t *testing.T) {
	handler, _, _, mockCommentUC, _ := setupTestHandler()
	
	tests := []struct {
		name           string
		userID         string
		articleID      string
		requestBody    interface{}
		mockSetup      func(*MockCommentUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "successful comment creation",
			userID:    "user123",
			articleID: "article123",
			requestBody: usecase.CreateCommentRequest{
				Content: "This is a test comment",
			},
			mockSetup: func(mockUC *MockCommentUseCase) {
				mockUC.On("Create", mock.Anything, "user123", "article123", mock.AnythingOfType("usecase.CreateCommentRequest")).
					Return(&usecase.CommentResponse{
						ID:        "comment123",
						ArticleID: "article123",
						AuthorID:  "user123",
						Content:   "This is a test comment",
						Status:    "pending",
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:      "invalid request body",
			userID:    "user123",
			articleID: "article123",
			requestBody: map[string]interface{}{
				"content": "", // invalid - too short
			},
			mockSetup:      func(mockUC *MockCommentUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockCommentUC.ExpectedCalls = nil
			tt.mockSetup(mockCommentUC)

			// Setup request
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/articles/"+tt.articleID+"/comments", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)
			c.SetParamNames("id")
			c.SetParamValues(tt.articleID)

			// Execute
			err := handler.CreateComment(c)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockCommentUC.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_RateArticle(t *testing.T) {
	handler, _, _, _, mockRatingUC := setupTestHandler()
	
	tests := []struct {
		name           string
		userID         string
		articleID      string
		requestBody    interface{}
		mockSetup      func(*MockRatingUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "successful rating",
			userID:    "user123",
			articleID: "article123",
			requestBody: usecase.RatingRequest{
				Stars: 5,
			},
			mockSetup: func(mockUC *MockRatingUseCase) {
				mockUC.On("RateArticle", mock.Anything, "user123", "article123", 5).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:      "invalid rating",
			userID:    "user123",
			articleID: "article123",
			requestBody: map[string]interface{}{
				"stars": 6, // invalid - too high
			},
			mockSetup:      func(mockUC *MockRatingUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockRatingUC.ExpectedCalls = nil
			tt.mockSetup(mockRatingUC)

			// Setup request
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/articles/"+tt.articleID+"/rate", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)
			c.SetParamNames("id")
			c.SetParamValues(tt.articleID)

			// Execute
			err := handler.RateArticle(c)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockRatingUC.AssertExpectations(t)
		})
	}
}
