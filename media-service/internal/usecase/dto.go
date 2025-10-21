package usecase

import "time"

// Upload
type UploadRequest struct {
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	Size         int64  `json:"size"`
	Data         []byte `json:"data"`
}

type MediaResponse struct {
	ID           string                 `json:"id"`
	UploaderID   string                 `json:"uploader_id"`
	Filename     string                 `json:"filename"`
	OriginalName string                 `json:"original_name"`
	MimeType     string                 `json:"mime_type"`
	Size         int64                  `json:"size"`
	Type         string                 `json:"type"`
	Status       string                 `json:"status"`
	URL          string                 `json:"url"`
	ThumbnailURL string                 `json:"thumbnail_url,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type UploadResponse struct {
	Media *MediaResponse `json:"media"`
	URL   string         `json:"url"`
}

// List
type ListFilter struct {
	UploaderID string
	Type       string
	Status     string
	Page       int
	PageSize   int
}

type ListResponse struct {
	Media []*MediaResponse `json:"media"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}
