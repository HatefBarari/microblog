package domain

import "time"

type ArticleStatus string

const (
	StatusDraft    ArticleStatus = "draft"
	StatusPending  ArticleStatus = "pending"
	StatusApproved ArticleStatus = "approved"
	StatusRejected ArticleStatus = "rejected"
	StatusArchived ArticleStatus = "archived"
)

type Article struct {
	ID          string        `bson:"_id,omitempty"`
	AuthorID    string        `bson:"author_id"`
	Title       string        `bson:"title"`
	Slug        string        `bson:"slug"`
	Summary     string        `bson:"summary"`
	Content     string        `bson:"content"`
	CoverURL    string        `bson:"cover_url"`
	Status      ArticleStatus `bson:"status"`
	CategoryID  string        `bson:"category_id"`
	Tags        []string      `bson:"tags"`
	ViewCount   int64         `bson:"view_count"`
	RatingAvg   float64       `bson:"rating_avg"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
	PublishedAt *time.Time    `bson:"published_at"`
}