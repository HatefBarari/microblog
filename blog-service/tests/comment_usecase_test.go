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

type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockCommentRepository) ListByArticle(ctx context.Context, articleID string, status domain.CommentStatus) ([]*domain.Comment, error) {
	args := m.Called(ctx, articleID, status)
	return args.Get(0).([]*domain.Comment), args.Error(1)
}

func (m *MockCommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *MockCommentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCommentUseCase_Create(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		userID         string
		articleID      string
		request        usecase.CreateCommentRequest
		mockSetup      func(*MockCommentRepository)
		expectedError  string
		expectedResult func(*usecase.CommentResponse) bool
	}{
		{
			name:      "successful comment creation",
			userID:    "user123",
			articleID: "article123",
			request: usecase.CreateCommentRequest{
				Content: "This is a test comment",
			},
			mockSetup: func(mockRepo *MockCommentRepository) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(comment *domain.Comment) bool {
					return comment.ArticleID == "article123" &&
						comment.AuthorID == "user123" &&
						comment.Content == "This is a test comment" &&
						comment.Status == domain.CommentPending
				})).Return(nil)
			},
			expectedError: "",
			expectedResult: func(resp *usecase.CommentResponse) bool {
				return resp.ArticleID == "article123" &&
					resp.AuthorID == "user123" &&
					resp.Content == "This is a test comment" &&
					resp.Status == string(domain.CommentPending)
			},
		},
		{
			name:      "successful reply creation",
			userID:    "user123",
			articleID: "article123",
			request: usecase.CreateCommentRequest{
				ParentID: "comment123",
				Content:  "This is a reply",
			},
			mockSetup: func(mockRepo *MockCommentRepository) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(comment *domain.Comment) bool {
					return comment.ArticleID == "article123" &&
						comment.AuthorID == "user123" &&
						comment.ParentID == "comment123" &&
						comment.Content == "This is a reply" &&
						comment.Status == domain.CommentPending
				})).Return(nil)
			},
			expectedError: "",
			expectedResult: func(resp *usecase.CommentResponse) bool {
				return resp.ArticleID == "article123" &&
					resp.AuthorID == "user123" &&
					resp.ParentID == "comment123" &&
					resp.Content == "This is a reply" &&
					resp.Status == string(domain.CommentPending)
			},
		},
		{
			name:      "repository error",
			userID:    "user123",
			articleID: "article123",
			request: usecase.CreateCommentRequest{
				Content: "This is a test comment",
			},
			mockSetup: func(mockRepo *MockCommentRepository) {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCommentRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewCommentUseCase(mockRepo, logger)

			result, err := uc.Create(context.Background(), tt.userID, tt.articleID, tt.request)

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

func TestCommentUseCase_ListByArticle(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		articleID      string
		status         domain.CommentStatus
		mockSetup      func(*MockCommentRepository)
		expectedError  string
		expectedResult func([]*usecase.CommentResponse) bool
	}{
		{
			name:      "successful list comments",
			articleID: "article123",
			status:    domain.CommentApproved,
			mockSetup: func(mockRepo *MockCommentRepository) {
				comments := []*domain.Comment{
					{
						ID:        "comment1",
						ArticleID: "article123",
						AuthorID:  "user1",
						Content:   "First comment",
						Status:    domain.CommentApproved,
						CreatedAt: time.Now(),
					},
					{
						ID:        "comment2",
						ArticleID: "article123",
						AuthorID:  "user2",
						Content:   "Second comment",
						Status:    domain.CommentApproved,
						CreatedAt: time.Now(),
					},
				}
				mockRepo.On("ListByArticle", mock.Anything, "article123", domain.CommentApproved).
					Return(comments, nil)
			},
			expectedError: "",
			expectedResult: func(resp []*usecase.CommentResponse) bool {
				return len(resp) == 2 &&
					resp[0].Content == "First comment" &&
					resp[1].Content == "Second comment"
			},
		},
		{
			name:      "empty comments list",
			articleID: "article123",
			status:    domain.CommentApproved,
			mockSetup: func(mockRepo *MockCommentRepository) {
				mockRepo.On("ListByArticle", mock.Anything, "article123", domain.CommentApproved).
					Return([]*domain.Comment{}, nil)
			},
			expectedError: "",
			expectedResult: func(resp []*usecase.CommentResponse) bool {
				return len(resp) == 0
			},
		},
		{
			name:      "repository error",
			articleID: "article123",
			status:    domain.CommentApproved,
			mockSetup: func(mockRepo *MockCommentRepository) {
				mockRepo.On("ListByArticle", mock.Anything, "article123", domain.CommentApproved).
					Return([]*domain.Comment{}, errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCommentRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewCommentUseCase(mockRepo, logger)

			result, err := uc.ListByArticle(context.Background(), tt.articleID, tt.status)

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

func TestCommentUseCase_UpdateStatus(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name          string
		commentID     string
		status        domain.CommentStatus
		mockSetup     func(*MockCommentRepository)
		expectedError string
	}{
		{
			name:      "successful status update",
			commentID: "comment123",
			status:    domain.CommentApproved,
			mockSetup: func(mockRepo *MockCommentRepository) {
				existingComment := &domain.Comment{
					ID:        "comment123",
					ArticleID: "article123",
					AuthorID:  "user123",
					Content:   "Test comment",
					Status:    domain.CommentPending,
				}
				mockRepo.On("GetByID", mock.Anything, "comment123").Return(existingComment, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(comment *domain.Comment) bool {
					return comment.ID == "comment123" &&
						comment.Status == domain.CommentApproved
				})).Return(nil)
			},
			expectedError: "",
		},
		{
			name:      "comment not found",
			commentID: "nonexistent",
			status:    domain.CommentApproved,
			mockSetup: func(mockRepo *MockCommentRepository) {
				mockRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedError: "comment not found",
		},
		{
			name:      "repository error on get",
			commentID: "comment123",
			status:    domain.CommentApproved,
			mockSetup: func(mockRepo *MockCommentRepository) {
				mockRepo.On("GetByID", mock.Anything, "comment123").Return(nil, errors.New("database error"))
			},
			expectedError: "database error",
		},
		{
			name:      "repository error on update",
			commentID: "comment123",
			status:    domain.CommentApproved,
			mockSetup: func(mockRepo *MockCommentRepository) {
				existingComment := &domain.Comment{
					ID:        "comment123",
					ArticleID: "article123",
					AuthorID:  "user123",
					Content:   "Test comment",
					Status:    domain.CommentPending,
				}
				mockRepo.On("GetByID", mock.Anything, "comment123").Return(existingComment, nil)
				mockRepo.On("Update", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCommentRepository)
			tt.mockSetup(mockRepo)

			uc := usecase.NewCommentUseCase(mockRepo, logger)

			err := uc.UpdateStatus(context.Background(), tt.commentID, tt.status)

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
