package domain

import "time"

type Rating struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"user_id"`
	TargetID  string    `bson:"target_id"` // article or comment
	Type      string    `bson:"type"`      // "article" or "comment"
	Stars     int       `bson:"stars"`     // 1..5
	CreatedAt time.Time `bson:"created_at"`
}