package domain

import "time"

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
	MediaTypeDocument MediaType = "document"
)

type MediaStatus string

const (
	MediaStatusPending   MediaStatus = "pending"
	MediaStatusProcessed MediaStatus = "processed"
	MediaStatusFailed    MediaStatus = "failed"
	MediaStatusDeleted   MediaStatus = "deleted"
)

type Media struct {
	ID          string     `bson:"_id,omitempty"`
	UploaderID  string     `bson:"uploader_id"`
	Filename    string     `bson:"filename"`
	OriginalName string    `bson:"original_name"`
	MimeType    string     `bson:"mime_type"`
	Size        int64      `bson:"size"`
	Type        MediaType  `bson:"type"`
	Status      MediaStatus `bson:"status"`
	URL         string     `bson:"url"`
	ThumbnailURL string    `bson:"thumbnail_url,omitempty"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty"`
	CreatedAt   time.Time  `bson:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at"`
}

type MediaRepository interface {
	Create(ctx context.Context, media *Media) error
	GetByID(ctx context.Context, id string) (*Media, error)
	GetByURL(ctx context.Context, url string) (*Media, error)
	ListByUploader(ctx context.Context, uploaderID string, limit, offset int) ([]*Media, int, error)
	Update(ctx context.Context, media *Media) error
	Delete(ctx context.Context, id string) error
	DeleteByUploader(ctx context.Context, uploaderID, id string) error
}
