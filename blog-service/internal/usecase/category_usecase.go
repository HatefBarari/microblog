package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/HatefBarari/microblog-blog/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type CategoryUseCase struct {
	repo domain.CategoryRepository
	logLevel      string
	// cfg  *Config
	log  *zap.Logger
}

func NewCategoryUseCase(repo domain.CategoryRepository, logLevel string, log *zap.Logger) *CategoryUseCase {
	return &CategoryUseCase{repo: repo, logLevel: logLevel, log: log}
}

func (uc *CategoryUseCase) Create(ctx context.Context, name, parentID string) (*CategoryResponse, error) {
	slug := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	// check duplicate slug
	exist, _ := uc.repo.GetBySlug(ctx, slug)
	if exist != nil {
		return nil, errors.New("slug already exists")
	}
	c := &domain.Category{
		ID:       primitive.NewObjectID().Hex(),
		Name:     name,
		Slug:     slug,
		ParentID: parentID,
	}
	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return &CategoryResponse{
		ID:       c.ID,
		Name:     c.Name,
		Slug:     c.Slug,
		ParentID: c.ParentID,
	}, nil
}

func (uc *CategoryUseCase) ListTree(ctx context.Context) ([]CategoryResponse, error) {
	list, err := uc.repo.ListTree(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]CategoryResponse, len(list))
	for i, c := range list {
		resp[i] = CategoryResponse{
			ID:       c.ID,
			Name:     c.Name,
			Slug:     c.Slug,
			ParentID: c.ParentID,
		}
	}
	return resp, nil
}