package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/HatefBarari/microblog-blog/internal/domain"
	"go.uber.org/zap"
)

type RatingUseCase struct {
	repo  domain.RatingRepository
	logLevel string
	log   *zap.Logger
}

func NewRatingUseCase(repo domain.RatingRepository, logLevel string, log *zap.Logger) *RatingUseCase {
	return &RatingUseCase{repo: repo, logLevel: logLevel, log: log}
}

func (uc *RatingUseCase) RateArticle(ctx context.Context, userID, articleID string, stars int) error {
	if stars < 1 || stars > 5 {
		return errors.New("stars must be between 1 and 5")
	}
	exist, _ := uc.repo.GetByUserAndTarget(ctx, userID, articleID, "article")
	if exist != nil {
		// آپدیت
		exist.Stars = stars
		return uc.repo.Save(ctx, exist)
	}
	// insert
	r := &domain.Rating{
		UserID:    userID,
		TargetID:  articleID,
		Type:      "article",
		Stars:     stars,
		CreatedAt: time.Now(),
	}
	return uc.repo.Save(ctx, r)
}

func (uc *RatingUseCase) Delete(ctx context.Context, userID, targetID, targetType string) error {
	return uc.repo.Delete(ctx, userID, targetID, targetType)
}