package usecase

import "time"

// Article
type CreateArticleRequest struct {
	Title      string   `json:"title" validate:"required,min=3"`
	Content    string   `json:"content" validate:"required,min=10"`
	Summary    string   `json:"summary" validate:"required,min=10"`
	CategoryID string   `json:"category_id" validate:"required"`
	Tags       []string `json:"tags"`
	CoverURL   string   `json:"cover_url"`
}

type ArticleResponse struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"author_id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Summary     string    `json:"summary"`
	Content     string    `json:"content"`
	CoverURL    string    `json:"cover_url"`
	Status      string    `json:"status"`
	CategoryID  string    `json:"category_id"`
	Tags        []string  `json:"tags"`
	ViewCount   int64     `json:"view_count"`
	RatingAvg   float64   `json:"rating_avg"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category
type CategoryResponse struct {
	ID       string                `json:"id"`
	Name     string                `json:"name"`
	Slug     string                `json:"slug"`
	ParentID string                `json:"parent_id,omitempty"`
	Children []CategoryResponse    `json:"children,omitempty"`
}

// Comment
type CreateCommentRequest struct {
	ParentID string `json:"parent_id,omitempty"` // فقط یک سطح
	Content  string `json:"content" validate:"required,min=2"`
}

type CommentResponse struct {
	ID        string    `json:"id"`
	ArticleID string    `json:"article_id"`
	ParentID  string    `json:"parent_id,omitempty"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Rating
type RatingRequest struct {
	Stars int `json:"stars" validate:"required,min=1,max=5"`
}
type ListFilter struct {
	AuthorID   *string
	CategoryID *string
	Status     *string
	Tag        *string
	Page       int
	PageSize   int
}