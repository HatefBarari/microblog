package domain

import "time"

type CommentStatus string

const (
	CommentPending  CommentStatus = "pending"
	CommentApproved CommentStatus = "approved"
	CommentRejected CommentStatus = "rejected"
)

type Comment struct {
	ID        string        `bson:"_id,omitempty"`
	ArticleID string        `bson:"article_id"`
	ParentID  string        `bson:"parent_id"` // فقط یک سطح
	AuthorID  string        `bson:"author_id"`
	Content   string        `bson:"content"`
	Status    CommentStatus `bson:"status"`
	CreatedAt time.Time     `bson:"created_at"`
}