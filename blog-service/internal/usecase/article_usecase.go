package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/HatefBarari/microblog-blog/internal/domain"
	"go.uber.org/zap"
)

type ArticleUseCase struct {
	repo          domain.ArticleRepository
	ratingRepo    domain.RatingRepository
	accessSecret  string
	refreshSecret string
	logLevel      string
	log           *zap.Logger
}

func NewArticleUseCase(repo domain.ArticleRepository, ratingRepo domain.RatingRepository, accessSecret, refreshSecret, logLevel string, log *zap.Logger) *ArticleUseCase {
	return &ArticleUseCase{
		repo:          repo,
		ratingRepo:    ratingRepo,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		logLevel:      logLevel,
		log:           log,
	}
}

func (uc *ArticleUseCase) Create(ctx context.Context, userID string, req CreateArticleRequest) (*ArticleResponse, error) {
	// ۱) Validation
	if len(req.Tags) > 10 {
		return nil, errors.New("max 10 tags")
	}
	// ۲) generate unique slug
	base := strings.ReplaceAll(strings.ToLower(req.Title), " ", "-")
	slug := base
	for i := 0; i < 5; i++ {
		exist, _ := uc.repo.GetBySlug(ctx, slug)
		if exist == nil {
			break
		}
		slug = fmt.Sprintf("%s-%d", base, i+1)
	}
	// ۳) build entity
	a := &domain.Article{
		AuthorID:   userID,
		Title:      req.Title,
		Slug:       slug,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverURL:   req.CoverURL,
		Status:     domain.StatusDraft,
		CategoryID: req.CategoryID,
		Tags:       req.Tags,
		ViewCount:  0,
		RatingAvg:  0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	// ۴) save
	if err := uc.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	// ۵) build response
	return &ArticleResponse{
		ID:        a.ID,
		AuthorID:  a.AuthorID,
		Title:     a.Title,
		Slug:      a.Slug,
		Summary:   a.Summary,
		Content:   a.Content,
		CoverURL:  a.CoverURL,
		Status:    string(a.Status),
		CategoryID: a.CategoryID,
		Tags:      a.Tags,
		ViewCount: a.ViewCount,
		RatingAvg: a.RatingAvg,
		CreatedAt: a.CreatedAt,
	}, nil
}

func (uc *ArticleUseCase) GetBySlug(ctx context.Context, slug string) (*ArticleResponse, error) {
	a, err := uc.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("article not found")
	}
	// bump view
	_ = uc.repo.UpdateViewCount(ctx, a.ID)
	return &ArticleResponse{
		ID:        a.ID,
		AuthorID:  a.AuthorID,
		Title:     a.Title,
		Slug:      a.Slug,
		Summary:   a.Summary,
		Content:   a.Content,
		CoverURL:  a.CoverURL,
		Status:    string(a.Status),
		CategoryID: a.CategoryID,
		Tags:      a.Tags,
		ViewCount: a.ViewCount + 1,
		RatingAvg: a.RatingAvg,
		CreatedAt: a.CreatedAt,
	}, nil
}

func (uc *ArticleUseCase) List(ctx context.Context, filter domain.ListFilter) ([]*ArticleResponse, int, error) {
	list, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]*ArticleResponse, len(list))
	for i, a := range list {
		resp[i] = &ArticleResponse{
			ID:        a.ID,
			AuthorID:  a.AuthorID,
			Title:     a.Title,
			Slug:      a.Slug,
			Summary:   a.Summary,
			Content:   a.Content,
			CoverURL:  a.CoverURL,
			Status:    string(a.Status),
			CategoryID: a.CategoryID,
			Tags:      a.Tags,
			ViewCount: a.ViewCount,
			RatingAvg: a.RatingAvg,
			CreatedAt: a.CreatedAt,
		}
	}
	return resp, total, nil
}

func (uc *ArticleUseCase) Update(ctx context.Context, userID string, id string, req CreateArticleRequest) (*ArticleResponse, error) {
	// چک مالکیت
	a, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, errors.New("article not found")
	}
	if a.AuthorID != userID {
		return nil, errors.New("forbidden")
	}
	// update fields
	a.Title = req.Title
	a.Content = req.Content
	a.Summary = req.Summary
	a.CategoryID = req.CategoryID
	a.Tags = req.Tags
	a.CoverURL = req.CoverURL
	a.UpdatedAt = time.Now()
	if err := uc.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return &ArticleResponse{
		ID:        a.ID,
		AuthorID:  a.AuthorID,
		Title:     a.Title,
		Slug:      a.Slug,
		Summary:   a.Summary,
		Content:   a.Content,
		CoverURL:  a.CoverURL,
		Status:    string(a.Status),
		CategoryID: a.CategoryID,
		Tags:      a.Tags,
		ViewCount: a.ViewCount,
		RatingAvg: a.RatingAvg,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}, nil
}

func (uc *ArticleUseCase) Delete(ctx context.Context, userID string, id string) error {
	a, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if a == nil {
		return errors.New("article not found")
	}
	if a.AuthorID != userID {
		return errors.New("forbidden")
	}
	return uc.repo.Delete(ctx, id)
}