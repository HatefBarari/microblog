package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/HatefBarari/microblog-media/internal/domain"
	"go.uber.org/zap"
)

type MediaUseCase struct {
	repo     domain.MediaRepository
	storage  MediaStorage
	cfg      *Config
	log      *zap.Logger
}

type Config struct {
	MaxFileSize    int64
	AllowedTypes   []string
	StoragePath    string
	BaseURL        string
	ThumbnailSize  int
}

type MediaStorage interface {
	Save(ctx context.Context, filename string, data []byte) (string, error)
	Get(ctx context.Context, filename string) ([]byte, error)
	Delete(ctx context.Context, filename string) error
	GenerateThumbnail(ctx context.Context, filename string, size int) (string, error)
}

func NewMediaUseCase(repo domain.MediaRepository, storage MediaStorage, cfg *Config, log *zap.Logger) *MediaUseCase {
	return &MediaUseCase{
		repo:    repo,
		storage: storage,
		cfg:     cfg,
		log:     log,
	}
}

func (uc *MediaUseCase) Upload(ctx context.Context, uploaderID string, file *multipart.FileHeader) (*UploadResponse, error) {
	// 1. Validate file
	if err := uc.validateFile(file); err != nil {
		return nil, err
	}

	// 2. Read file data
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 3. Generate unique filename
	filename := uc.generateFilename(file.Filename)
	
	// 4. Save to storage
	url, err := uc.storage.Save(ctx, filename, data)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 5. Generate thumbnail for images
	var thumbnailURL string
	if uc.isImage(file.Header.Get("Content-Type")) {
		thumbURL, err := uc.storage.GenerateThumbnail(ctx, filename, uc.cfg.ThumbnailSize)
		if err != nil {
			uc.log.Warn("failed to generate thumbnail", zap.Error(err))
		} else {
			thumbnailURL = thumbURL
		}
	}

	// 6. Create media record
	media := &domain.Media{
		UploaderID:  uploaderID,
		Filename:    filename,
		OriginalName: file.Filename,
		MimeType:    file.Header.Get("Content-Type"),
		Size:        file.Size,
		Type:        uc.getMediaType(file.Header.Get("Content-Type")),
		Status:      domain.MediaStatusProcessed,
		URL:         url,
		ThumbnailURL: thumbnailURL,
		Metadata: map[string]interface{}{
			"upload_ip": ctx.Value("client_ip"),
			"user_agent": ctx.Value("user_agent"),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, media); err != nil {
		// Clean up uploaded file if database save fails
		_ = uc.storage.Delete(ctx, filename)
		return nil, fmt.Errorf("failed to save media record: %w", err)
	}

	// 7. Build response
	response := &MediaResponse{
		ID:           media.ID,
		UploaderID:   media.UploaderID,
		Filename:     media.Filename,
		OriginalName: media.OriginalName,
		MimeType:     media.MimeType,
		Size:         media.Size,
		Type:         string(media.Type),
		Status:       string(media.Status),
		URL:          media.URL,
		ThumbnailURL: media.ThumbnailURL,
		Metadata:     media.Metadata,
		CreatedAt:    media.CreatedAt,
		UpdatedAt:    media.UpdatedAt,
	}

	return &UploadResponse{
		Media: response,
		URL:   url,
	}, nil
}

func (uc *MediaUseCase) GetByID(ctx context.Context, id string) (*MediaResponse, error) {
	media, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if media == nil {
		return nil, errors.New("media not found")
	}

	return &MediaResponse{
		ID:           media.ID,
		UploaderID:   media.UploaderID,
		Filename:     media.Filename,
		OriginalName: media.OriginalName,
		MimeType:     media.MimeType,
		Size:         media.Size,
		Type:         string(media.Type),
		Status:       string(media.Status),
		URL:          media.URL,
		ThumbnailURL: media.ThumbnailURL,
		Metadata:     media.Metadata,
		CreatedAt:    media.CreatedAt,
		UpdatedAt:    media.UpdatedAt,
	}, nil
}

func (uc *MediaUseCase) List(ctx context.Context, filter ListFilter) (*ListResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	offset := (filter.Page - 1) * filter.PageSize
	media, total, err := uc.repo.ListByUploader(ctx, filter.UploaderID, filter.PageSize, offset)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := make([]*MediaResponse, len(media))
	for i, m := range media {
		response[i] = &MediaResponse{
			ID:           m.ID,
			UploaderID:   m.UploaderID,
			Filename:     m.Filename,
			OriginalName: m.OriginalName,
			MimeType:     m.MimeType,
			Size:         m.Size,
			Type:         string(m.Type),
			Status:       string(m.Status),
			URL:          m.URL,
			ThumbnailURL: m.ThumbnailURL,
			Metadata:     m.Metadata,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		}
	}

	return &ListResponse{
		Media: response,
		Total: total,
		Page:  filter.Page,
		Size:  len(response),
	}, nil
}

func (uc *MediaUseCase) Delete(ctx context.Context, uploaderID, id string) error {
	// Check ownership
	media, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if media == nil {
		return errors.New("media not found")
	}
	if media.UploaderID != uploaderID {
		return errors.New("forbidden")
	}

	// Delete from storage
	if err := uc.storage.Delete(ctx, media.Filename); err != nil {
		uc.log.Warn("failed to delete file from storage", zap.Error(err))
	}

	// Delete from database
	return uc.repo.DeleteByUploader(ctx, uploaderID, id)
}

func (uc *MediaUseCase) validateFile(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > uc.cfg.MaxFileSize {
		return fmt.Errorf("file too large: %d bytes (max: %d)", file.Size, uc.cfg.MaxFileSize)
	}

	// Check file type
	contentType := file.Header.Get("Content-Type")
	if !uc.isAllowedType(contentType) {
		return fmt.Errorf("file type not allowed: %s", contentType)
	}

	return nil
}

func (uc *MediaUseCase) isAllowedType(contentType string) bool {
	for _, allowedType := range uc.cfg.AllowedTypes {
		if strings.HasPrefix(contentType, allowedType) {
			return true
		}
	}
	return false
}

func (uc *MediaUseCase) isImage(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

func (uc *MediaUseCase) getMediaType(contentType string) domain.MediaType {
	if strings.HasPrefix(contentType, "image/") {
		return domain.MediaTypeImage
	}
	if strings.HasPrefix(contentType, "video/") {
		return domain.MediaTypeVideo
	}
	if strings.HasPrefix(contentType, "audio/") {
		return domain.MediaTypeAudio
	}
	return domain.MediaTypeDocument
}

func (uc *MediaUseCase) generateFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d_%s%s", timestamp, strings.ReplaceAll(originalName, " ", "_"), ext)
}
