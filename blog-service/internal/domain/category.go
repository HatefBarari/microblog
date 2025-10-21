package domain

type Category struct {
	ID       string `bson:"_id,omitempty"`
	Name     string `bson:"name"`
	Slug     string `bson:"slug"`
	ParentID string `bson:"parent_id"` // درختی
}