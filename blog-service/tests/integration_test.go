package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HatefBarari/microblog-blog/internal/domain"
	"github.com/HatefBarari/microblog-blog/internal/presenter"
	"github.com/HatefBarari/microblog-blog/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Integration test for the complete blog service flow
func TestBlogServiceIntegration(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	// Setup mocks
	mockArticleRepo := new(MockArticleRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockCommentRepo := new(MockCommentRepository)
	mockRatingRepo := new(MockRatingRepository)
	
	// Setup use cases
	articleUC := usecase.NewArticleUseCase(mockArticleRepo, mockRatingRepo, "access_secret", "refresh_secret", "debug", logger)
	categoryUC := usecase.NewCategoryUseCase(mockCategoryRepo, logger)
	commentUC := usecase.NewCommentUseCase(mockCommentRepo, logger)
	ratingUC := usecase.NewRatingUseCase(mockRatingRepo, logger)
	
	// Setup handler
	handler := presenter.NewHTTPHandler(articleUC, categoryUC, commentUC, ratingUC)
	
	// Setup echo
	e := echo.New()
	
	t.Run("Complete Article Lifecycle", func(t *testing.T) {
		// 1. Create article
		createReq := usecase.CreateArticleRequest{
			Title:      "Integration Test Article",
			Content:    "This is a comprehensive test article for integration testing",
			Summary:    "Test summary for integration",
			CategoryID: "cat123",
			Tags:       []string{"test", "integration", "go"},
			CoverURL:   "https://example.com/cover.jpg",
		}
		
		// Mock article creation
		mockArticleRepo.On("GetBySlug", mock.Anything, "integration-test-article").Return(nil, nil)
		mockArticleRepo.On("Create", mock.Anything, mock.MatchedBy(func(article *domain.Article) bool {
			return article.Title == "Integration Test Article" &&
				article.Slug == "integration-test-article" &&
				article.AuthorID == "user123"
		})).Return(nil)
		
		// Create request
		reqBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		
		c := e.NewContext(req, rec)
		c.Set("userID", "user123")
		
		err := handler.CreateArticle(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		
		// 2. Get article by slug
		mockArticleRepo.On("GetBySlug", mock.Anything, "integration-test-article").Return(&domain.Article{
			ID:        "article123",
			AuthorID:  "user123",
			Title:     "Integration Test Article",
			Slug:      "integration-test-article",
			Summary:   "Test summary for integration",
			Content:   "This is a comprehensive test article for integration testing",
			CoverURL:  "https://example.com/cover.jpg",
			Status:    domain.StatusDraft,
			ViewCount: 0,
			RatingAvg: 0,
		}, nil)
		mockArticleRepo.On("UpdateViewCount", mock.Anything, "article123").Return(nil)
		
		req = httptest.NewRequest(http.MethodGet, "/articles/integration-test-article", nil)
		rec = httptest.NewRecorder()
		
		c = e.NewContext(req, rec)
		c.SetParamNames("slug")
		c.SetParamValues("integration-test-article")
		
		err = handler.GetArticleBySlug(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		
		// 3. Add comment
		commentReq := usecase.CreateCommentRequest{
			Content: "Great article!",
		}
		
		mockCommentRepo.On("Create", mock.Anything, mock.MatchedBy(func(comment *domain.Comment) bool {
			return comment.ArticleID == "article123" &&
				comment.AuthorID == "user123" &&
				comment.Content == "Great article!"
		})).Return(nil)
		
		reqBody, _ = json.Marshal(commentReq)
		req = httptest.NewRequest(http.MethodPost, "/articles/article123/comments", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		c = e.NewContext(req, rec)
		c.Set("userID", "user123")
		c.SetParamNames("id")
		c.SetParamValues("article123")
		
		err = handler.CreateComment(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		
		// 4. Rate article
		ratingReq := usecase.RatingRequest{
			Stars: 5,
		}
		
		mockRatingRepo.On("GetByUserAndTarget", mock.Anything, "user123", "article123", "article").
			Return(nil, nil)
		mockRatingRepo.On("Create", mock.Anything, mock.MatchedBy(func(rating *domain.Rating) bool {
			return rating.UserID == "user123" &&
				rating.TargetID == "article123" &&
				rating.Type == "article" &&
				rating.Stars == 5
		})).Return(nil)
		mockRatingRepo.On("GetAverageByTarget", mock.Anything, "article123", "article").
			Return(5.0, nil)
		mockRatingRepo.On("UpdateTargetAverage", mock.Anything, "article123", "article", 5.0).
			Return(nil)
		
		reqBody, _ = json.Marshal(ratingReq)
		req = httptest.NewRequest(http.MethodPost, "/articles/article123/rate", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		c = e.NewContext(req, rec)
		c.Set("userID", "user123")
		c.SetParamNames("id")
		c.SetParamValues("article123")
		
		err = handler.RateArticle(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		
		// 5. Update article
		updateReq := usecase.CreateArticleRequest{
			Title:      "Updated Integration Test Article",
			Content:    "This is an updated comprehensive test article",
			Summary:    "Updated test summary for integration",
			CategoryID: "cat456",
			Tags:       []string{"test", "integration", "go", "updated"},
		}
		
		mockArticleRepo.On("GetByID", mock.Anything, "article123").Return(&domain.Article{
			ID:       "article123",
			AuthorID: "user123",
			Title:    "Integration Test Article",
		}, nil)
		mockArticleRepo.On("Update", mock.Anything, mock.MatchedBy(func(article *domain.Article) bool {
			return article.ID == "article123" &&
				article.Title == "Updated Integration Test Article"
		})).Return(nil)
		
		reqBody, _ = json.Marshal(updateReq)
		req = httptest.NewRequest(http.MethodPut, "/articles/article123", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		c = e.NewContext(req, rec)
		c.Set("userID", "user123")
		c.SetParamNames("id")
		c.SetParamValues("article123")
		
		err = handler.UpdateArticle(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		
		// 6. List articles
		mockArticleRepo.On("List", mock.Anything, mock.AnythingOfType("domain.ListFilter")).
			Return([]*domain.Article{
				{
					ID:        "article123",
					AuthorID:  "user123",
					Title:     "Updated Integration Test Article",
					Slug:      "integration-test-article",
					Status:    domain.StatusDraft,
					ViewCount: 1,
					RatingAvg: 5.0,
				},
			}, 1, nil)
		
		req = httptest.NewRequest(http.MethodGet, "/articles", nil)
		rec = httptest.NewRecorder()
		
		c = e.NewContext(req, rec)
		
		err = handler.ListArticles(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		
		// 7. Delete article
		mockArticleRepo.On("GetByID", mock.Anything, "article123").Return(&domain.Article{
			ID:       "article123",
			AuthorID: "user123",
		}, nil)
		mockArticleRepo.On("Delete", mock.Anything, "article123").Return(nil)
		
		req = httptest.NewRequest(http.MethodDelete, "/articles/article123", nil)
		rec = httptest.NewRecorder()
		
		c = e.NewContext(req, rec)
		c.Set("userID", "user123")
		c.SetParamNames("id")
		c.SetParamValues("article123")
		
		err = handler.DeleteArticle(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		
		// Assert all expectations
		mockArticleRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
		mockRatingRepo.AssertExpectations(t)
	})
	
	t.Run("Category Management Flow", func(t *testing.T) {
		// 1. Create root category
		mockCategoryRepo.On("GetBySlug", mock.Anything, "technology").Return(nil, nil)
		mockCategoryRepo.On("Create", mock.Anything, mock.MatchedBy(func(category *domain.Category) bool {
			return category.Name == "Technology" &&
				category.Slug == "technology" &&
				category.ParentID == ""
		})).Return(nil)
		
		reqBody, _ := json.Marshal(map[string]interface{}{
			"name": "Technology",
		})
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		
		c := e.NewContext(req, rec)
		c.Set("role", "admin")
		
		err := handler.CreateCategory(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		
		// 2. Create subcategory
		mockCategoryRepo.On("GetBySlug", mock.Anything, "programming").Return(nil, nil)
		mockCategoryRepo.On("Create", mock.Anything, mock.MatchedBy(func(category *domain.Category) bool {
			return category.Name == "Programming" &&
				category.Slug == "programming" &&
				category.ParentID == "cat123"
		})).Return(nil)
		
		reqBody, _ = json.Marshal(map[string]interface{}{
			"name":      "Programming",
			"parent_id": "cat123",
		})
		req = httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		c = e.NewContext(req, rec)
		c.Set("role", "admin")
		
		err = handler.CreateCategory(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		
		// 3. List category tree
		mockCategoryRepo.On("List", mock.Anything).Return([]*domain.Category{
			{ID: "cat1", Name: "Technology", Slug: "technology", ParentID: ""},
			{ID: "cat2", Name: "Programming", Slug: "programming", ParentID: "cat1"},
			{ID: "cat3", Name: "Science", Slug: "science", ParentID: ""},
		}, nil)
		
		req = httptest.NewRequest(http.MethodGet, "/categories/tree", nil)
		rec = httptest.NewRecorder()
		
		c = e.NewContext(req, rec)
		
		err = handler.ListCategoryTree(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		
		mockCategoryRepo.AssertExpectations(t)
	})
	
	t.Run("Comment Management Flow", func(t *testing.T) {
		// 1. List comments
		mockCommentRepo.On("ListByArticle", mock.Anything, "article123", domain.CommentApproved).
			Return([]*domain.Comment{
				{
					ID:        "comment1",
					ArticleID: "article123",
					AuthorID:  "user1",
					Content:   "Great article!",
					Status:    domain.CommentApproved,
				},
				{
					ID:        "comment2",
					ArticleID: "article123",
					AuthorID:  "user2",
					Content:   "Very informative",
					Status:    domain.CommentApproved,
				},
			}, nil)
		
		req := httptest.NewRequest(http.MethodGet, "/articles/article123/comments", nil)
		rec := httptest.NewRecorder()
		
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("article123")
		
		err := handler.ListComments(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		
		// 2. Update comment status
		mockCommentRepo.On("GetByID", mock.Anything, "comment1").Return(&domain.Comment{
			ID:        "comment1",
			ArticleID: "article123",
			AuthorID:  "user1",
			Content:   "Great article!",
			Status:    domain.CommentPending,
		}, nil)
		mockCommentRepo.On("Update", mock.Anything, mock.MatchedBy(func(comment *domain.Comment) bool {
			return comment.ID == "comment1" &&
				comment.Status == domain.CommentApproved
		})).Return(nil)
		
		reqBody, _ := json.Marshal(map[string]interface{}{
			"status": "approved",
		})
		req = httptest.NewRequest(http.MethodPut, "/comments/comment1/status", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		
		c = e.NewContext(req, rec)
		c.Set("role", "admin")
		c.SetParamNames("id")
		c.SetParamValues("comment1")
		
		err = handler.UpdateCommentStatus(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		
		mockCommentRepo.AssertExpectations(t)
	})
}
