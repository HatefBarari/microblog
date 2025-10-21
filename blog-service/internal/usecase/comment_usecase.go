package usecase

import (
	"context"
	"time"

	"github.com/HatefBarari/microblog-blog/internal/domain"
	"go.uber.org/zap"
)

type CommentUseCase struct {
	repo  domain.CommentRepository
	logLevel string
	log   *zap.Logger
}

func NewCommentUseCase(repo domain.CommentRepository, logLevel string, log *zap.Logger) *CommentUseCase {
	return &CommentUseCase{repo: repo, logLevel: logLevel, log: log}
}

func (uc *CommentUseCase) Create(ctx context.Context, userID, articleID string, req CreateCommentRequest) (*CommentResponse, error) {
	c := &domain.Comment{
		ArticleID: articleID,
		ParentID:  req.ParentID,
		AuthorID:  userID,
		Content:   req.Content,
		Status:    domain.CommentPending, // منتظر تایید
		CreatedAt: time.Now(),
	}
	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return &CommentResponse{
		ID:        c.ID,
		ArticleID: c.ArticleID,
		ParentID:  c.ParentID,
		AuthorID:  c.AuthorID,
		Content:   c.Content,
		Status:    string(c.Status),
		CreatedAt: c.CreatedAt,
	}, nil
}

func (uc *CommentUseCase) ListByArticle(ctx context.Context, articleID string, status domain.CommentStatus) ([]*CommentResponse, error) {
	list, err := uc.repo.ListByArticle(ctx, articleID, status)
	if err != nil {
		return nil, err
	}
	resp := make([]*CommentResponse, len(list))
	for i, c := range list {
		resp[i] = &CommentResponse{
			ID:        c.ID,
			ArticleID: c.ArticleID,
			ParentID:  c.ParentID,
			AuthorID:  c.AuthorID,
			Content:   c.Content,
			Status:    string(c.Status),
			CreatedAt: c.CreatedAt,
		}
	}
	return resp, nil
}

func (uc *CommentUseCase) UpdateStatus(ctx context.Context, id string, status domain.CommentStatus) error {
	return uc.repo.UpdateStatus(ctx, id, status)
}

func (uc *CommentUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}