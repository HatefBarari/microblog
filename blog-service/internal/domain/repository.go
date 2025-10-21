package domain

import "context"

type ArticleRepository interface {
	Create(ctx context.Context, a *Article) error
	GetByID(ctx context.Context, id string) (*Article, error)
	GetBySlug(ctx context.Context, slug string) (*Article, error)
	List(ctx context.Context, filter ListFilter) ([]*Article, int, error)
	Update(ctx context.Context, a *Article) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status ArticleStatus) error
	UpdateViewCount(ctx context.Context, id string) error
}

type CategoryRepository interface {
	Create(ctx context.Context, c *Category) error
	GetByID(ctx context.Context, id string) (*Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	ListTree(ctx context.Context) ([]*Category, error)
}

type CommentRepository interface {
	Create(ctx context.Context, c *Comment) error
	GetByID(ctx context.Context, id string) (*Comment, error)
	ListByArticle(ctx context.Context, articleID string, status CommentStatus) ([]*Comment, error)
	UpdateStatus(ctx context.Context, id string, status CommentStatus) error
	Delete(ctx context.Context, id string) error
}

type RatingRepository interface {
	Save(ctx context.Context, r *Rating) error
	GetByUserAndTarget(ctx context.Context, userID, targetID, targetType string) (*Rating, error)
	GetAverage(ctx context.Context, targetID, targetType string) (float64, error)
	Delete(ctx context.Context, userID, targetID, targetType string) error
}

type ListFilter struct {
	AuthorID   *string
	CategoryID *string
	Status     *string
	Tag        *string
	Page       int
	PageSize   int
}