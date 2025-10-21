package infrastructure

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/HatefBarari/microblog-media/internal/usecase"
	"go.uber.org/zap"
)

type FileStorage struct {
	basePath string
	baseURL  string
	log      *zap.Logger
}

func NewFileStorage(basePath, baseURL string, log *zap.Logger) *FileStorage {
	return &FileStorage{
		basePath: basePath,
		baseURL:  baseURL,
		log:      log,
	}
}

func (fs *FileStorage) Save(ctx context.Context, filename string, data []byte) (string, error) {
	// Create directory if not exists
	dir := filepath.Join(fs.basePath, "uploads")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Create file path
	filePath := filepath.Join(dir, filename)
	
	// Write file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Generate URL
	url := fmt.Sprintf("%s/uploads/%s", fs.baseURL, filename)
	
	fs.log.Info("file saved successfully", 
		zap.String("filename", filename),
		zap.String("url", url),
		zap.Int64("size", int64(len(data))))

	return url, nil
}

func (fs *FileStorage) Get(ctx context.Context, filename string) ([]byte, error) {
	filePath := filepath.Join(fs.basePath, "uploads", filename)
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", filename)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

func (fs *FileStorage) Delete(ctx context.Context, filename string) error {
	filePath := filepath.Join(fs.basePath, "uploads", filename)
	
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			fs.log.Warn("file not found for deletion", zap.String("filename", filename))
			return nil // Not an error if file doesn't exist
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	fs.log.Info("file deleted successfully", zap.String("filename", filename))
	return nil
}

func (fs *FileStorage) GenerateThumbnail(ctx context.Context, filename string, size int) (string, error) {
	// For now, return the original URL as thumbnail
	// In a real implementation, you would use an image processing library
	// like github.com/disintegration/imaging to create thumbnails
	
	// Check if it's an image file
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}
	
	for _, imgExt := range imageExts {
		if ext == imgExt {
			// For now, just return the original URL
			// In production, you would generate an actual thumbnail
			url := fmt.Sprintf("%s/uploads/%s", fs.baseURL, filename)
			return url, nil
		}
	}
	
	return "", fmt.Errorf("thumbnail generation not supported for file type: %s", ext)
}

// Implement usecase.MediaStorage interface
var _ usecase.MediaStorage = (*FileStorage)(nil)
