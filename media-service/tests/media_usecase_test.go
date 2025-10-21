package tests

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/HatefBarari/microblog-media/internal/domain"
	"github.com/HatefBarari/microblog-media/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock repositories and storage
type MockMediaRepository struct {
	mock.Mock
}

func (m *MockMediaRepository) Create(ctx context.Context, media *domain.Media) error {
	args := m.Called(ctx, media)
	return args.Error(0)
}

func (m *MockMediaRepository) GetByID(ctx context.Context, id string) (*domain.Media, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Media), args.Error(1)
}

func (m *MockMediaRepository) GetByURL(ctx context.Context, url string) (*domain.Media, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Media), args.Error(1)
}

func (m *MockMediaRepository) ListByUploader(ctx context.Context, uploaderID string, limit, offset int) ([]*domain.Media, int, error) {
	args := m.Called(ctx, uploaderID, limit, offset)
	return args.Get(0).([]*domain.Media), args.Int(1), args.Error(2)
}

func (m *MockMediaRepository) Update(ctx context.Context, media *domain.Media) error {
	args := m.Called(ctx, media)
	return args.Error(0)
}

func (m *MockMediaRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMediaRepository) DeleteByUploader(ctx context.Context, uploaderID, id string) error {
	args := m.Called(ctx, uploaderID, id)
	return args.Error(0)
}

type MockMediaStorage struct {
	mock.Mock
}

func (m *MockMediaStorage) Save(ctx context.Context, filename string, data []byte) (string, error) {
	args := m.Called(ctx, filename, data)
	return args.String(0), args.Error(1)
}

func (m *MockMediaStorage) Get(ctx context.Context, filename string) ([]byte, error) {
	args := m.Called(ctx, filename)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockMediaStorage) Delete(ctx context.Context, filename string) error {
	args := m.Called(ctx, filename)
	return args.Error(0)
}

func (m *MockMediaStorage) GenerateThumbnail(ctx context.Context, filename string, size int) (string, error) {
	args := m.Called(ctx, filename, size)
	return args.String(0), args.Error(1)
}

func TestMediaUseCase_Upload(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		uploaderID     string
		file           *multipart.FileHeader
		mockSetup      func(*MockMediaRepository, *MockMediaStorage)
		expectedError  string
		expectedResult func(*usecase.UploadResponse) bool
	}{
		{
			name:       "successful image upload",
			uploaderID: "user123",
			file: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024,
				Header:   map[string][]string{"Content-Type": {"image/jpeg"}},
			},
			mockSetup: func(mockRepo *MockMediaRepository, mockStorage *MockMediaStorage) {
				mockStorage.On("Save", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
					Return("http://localhost:8083/uploads/test.jpg", nil)
				mockStorage.On("GenerateThumbnail", mock.Anything, mock.AnythingOfType("string"), 300).
					Return("http://localhost:8083/uploads/thumb_test.jpg", nil)
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(media *domain.Media) bool {
					return media.UploaderID == "user123" &&
						media.OriginalName == "test.jpg" &&
						media.Type == domain.MediaTypeImage &&
						media.Status == domain.MediaStatusProcessed
				})).Return(nil)
			},
			expectedError: "",
			expectedResult: func(resp *usecase.UploadResponse) bool {
				return resp.Media.UploaderID == "user123" &&
					resp.Media.OriginalName == "test.jpg" &&
					resp.Media.Type == string(domain.MediaTypeImage) &&
					resp.URL == "http://localhost:8083/uploads/test.jpg"
			},
		},
		{
			name:       "file too large",
			uploaderID: "user123",
			file: &multipart.FileHeader{
				Filename: "large.jpg",
				Size:     20 * 1024 * 1024, // 20MB
				Header:   map[string][]string{"Content-Type": {"image/jpeg"}},
			},
			mockSetup: func(mockRepo *MockMediaRepository, mockStorage *MockMediaStorage) {
				// No repository calls expected for invalid input
			},
			expectedError: "file too large",
		},
		{
			name:       "unsupported file type",
			uploaderID: "user123",
			file: &multipart.FileHeader{
				Filename: "test.txt",
				Size:     1024,
				Header:   map[string][]string{"Content-Type": {"text/plain"}},
			},
			mockSetup: func(mockRepo *MockMediaRepository, mockStorage *MockMediaStorage) {
				// No repository calls expected for invalid input
			},
			expectedError: "file type not allowed",
		},
		{
			name:       "storage error",
			uploaderID: "user123",
			file: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024,
				Header:   map[string][]string{"Content-Type": {"image/jpeg"}},
			},
			mockSetup: func(mockRepo *MockMediaRepository, mockStorage *MockMediaStorage) {
				mockStorage.On("Save", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
					Return("", errors.New("storage error"))
			},
			expectedError: "failed to save file",
		},
		{
			name:       "database error",
			uploaderID: "user123",
			file: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024,
				Header:   map[string][]string{"Content-Type": {"image/jpeg"}},
			},
			mockSetup: func(mockRepo *MockMediaRepository, mockStorage *MockMediaStorage) {
				mockStorage.On("Save", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
					Return("http://localhost:8083/uploads/test.jpg", nil)
				mockStorage.On("GenerateThumbnail", mock.Anything, mock.AnythingOfType("string"), 300).
					Return("http://localhost:8083/uploads/thumb_test.jpg", nil)
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))
				mockStorage.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)
			},
			expectedError: "failed to save media record",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockMediaRepository)
			mockStorage := new(MockMediaStorage)
			tt.mockSetup(mockRepo, mockStorage)

			uc := usecase.NewMediaUseCase(
				mockRepo,
				mockStorage,
				&usecase.Config{
					MaxFileSize:   10 * 1024 * 1024, // 10MB
					AllowedTypes:  []string{"image/"},
					StoragePath:   "./uploads",
					BaseURL:       "http://localhost:8083",
					ThumbnailSize: 300,
				},
				logger,
			)

			result, err := uc.Upload(context.Background(), tt.uploaderID, tt.file)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.True(t, tt.expectedResult(result))
			}

			mockRepo.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestMediaUseCase_GetByID(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		mediaID        string
		mockSetup      func(*MockMediaRepository)
		expectedError  string
		expectedResult func(*usecase.MediaResponse) bool
	}{
		{
			name:    "successful get by ID",
			mediaID: "media123",
			mockSetup: func(mockRepo *MockMediaRepository) {
				media := &domain.Media{
					ID:           "media123",
					UploaderID:   "user123",
					Filename:     "test.jpg",
					OriginalName: "test.jpg",
					MimeType:     "image/jpeg",
					Size:         1024,
					Type:         domain.MediaTypeImage,
					Status:       domain.MediaStatusProcessed,
					URL:          "http://localhost:8083/uploads/test.jpg",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				mockRepo.On("GetByID", mock.Anything, "media123").Return(media, nil)
			},
			expectedError: "",
			expectedResult: func(resp *usecase.MediaResponse) bool {
				return resp.ID == "media123" &&
					resp.UploaderID == "user123" &&
					resp.Filename == "test.jpg" &&
					resp.Type == string(domain.MediaTypeImage)
			},
		},
		{
			name:    "media not found",
			mediaID: "nonexistent",
			mockSetup: func(mockRepo *MockMediaRepository) {
				mockRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedError: "media not found",
		},
		{
			name:    "database error",
			mediaID: "media123",
			mockSetup: func(mockRepo *MockMediaRepository) {
				mockRepo.On("GetByID", mock.Anything, "media123").Return(nil, errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockMediaRepository)
			mockStorage := new(MockMediaStorage)
			tt.mockSetup(mockRepo)

			uc := usecase.NewMediaUseCase(
				mockRepo,
				mockStorage,
				&usecase.Config{},
				logger,
			)

			result, err := uc.GetByID(context.Background(), tt.mediaID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.True(t, tt.expectedResult(result))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMediaUseCase_List(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name           string
		filter         usecase.ListFilter
		mockSetup      func(*MockMediaRepository)
		expectedError  string
		expectedResult func(*usecase.ListResponse) bool
	}{
		{
			name: "successful list",
			filter: usecase.ListFilter{
				UploaderID: "user123",
				Page:       1,
				PageSize:   10,
			},
			mockSetup: func(mockRepo *MockMediaRepository) {
				media := []*domain.Media{
					{
						ID:           "media1",
						UploaderID:   "user123",
						Filename:     "test1.jpg",
						OriginalName: "test1.jpg",
						Type:         domain.MediaTypeImage,
						Status:       domain.MediaStatusProcessed,
						CreatedAt:    time.Now(),
					},
					{
						ID:           "media2",
						UploaderID:   "user123",
						Filename:     "test2.jpg",
						OriginalName: "test2.jpg",
						Type:         domain.MediaTypeImage,
						Status:       domain.MediaStatusProcessed,
						CreatedAt:    time.Now(),
					},
				}
				mockRepo.On("ListByUploader", mock.Anything, "user123", 10, 0).Return(media, 2, nil)
			},
			expectedError: "",
			expectedResult: func(resp *usecase.ListResponse) bool {
				return len(resp.Media) == 2 &&
					resp.Total == 2 &&
					resp.Page == 1 &&
					resp.Size == 2
			},
		},
		{
			name: "empty list",
			filter: usecase.ListFilter{
				UploaderID: "user123",
				Page:       1,
				PageSize:   10,
			},
			mockSetup: func(mockRepo *MockMediaRepository) {
				mockRepo.On("ListByUploader", mock.Anything, "user123", 10, 0).Return([]*domain.Media{}, 0, nil)
			},
			expectedError: "",
			expectedResult: func(resp *usecase.ListResponse) bool {
				return len(resp.Media) == 0 &&
					resp.Total == 0
			},
		},
		{
			name: "database error",
			filter: usecase.ListFilter{
				UploaderID: "user123",
				Page:       1,
				PageSize:   10,
			},
			mockSetup: func(mockRepo *MockMediaRepository) {
				mockRepo.On("ListByUploader", mock.Anything, "user123", 10, 0).Return([]*domain.Media{}, 0, errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockMediaRepository)
			mockStorage := new(MockMediaStorage)
			tt.mockSetup(mockRepo)

			uc := usecase.NewMediaUseCase(
				mockRepo,
				mockStorage,
				&usecase.Config{},
				logger,
			)

			result, err := uc.List(context.Background(), tt.filter)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.True(t, tt.expectedResult(result))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMediaUseCase_Delete(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	
	tests := []struct {
		name          string
		uploaderID    string
		mediaID       string
		mockSetup     func(*MockMediaRepository, *MockMediaStorage)
		expectedError string
	}{
		{
			name:       "successful delete",
			uploaderID: "user123",
			mediaID:    "media123",
			mockSetup: func(mockRepo *MockMediaRepository, mockStorage *MockMediaStorage) {
				media := &domain.Media{
					ID:         "media123",
					UploaderID: "user123",
					Filename:   "test.jpg",
				}
				mockRepo.On("GetByID", mock.Anything, "media123").Return(media, nil)
				mockStorage.On("Delete", mock.Anything, "test.jpg").Return(nil)
				mockRepo.On("DeleteByUploader", mock.Anything, "user123", "media123").Return(nil)
			},
			expectedError: "",
		},
		{
			name:       "media not found",
			uploaderID: "user123",
			mediaID:    "nonexistent",
			mockSetup: func(mockRepo *MockMediaRepository, mockStorage *MockMediaStorage) {
				mockRepo.On("GetByID", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectedError: "media not found",
		},
		{
			name:       "forbidden - not owner",
			uploaderID: "user456",
			mediaID:    "media123",
			mockSetup: func(mockRepo *MockMediaRepository, mockStorage *MockMediaStorage) {
				media := &domain.Media{
					ID:         "media123",
					UploaderID: "user123", // different owner
					Filename:   "test.jpg",
				}
				mockRepo.On("GetByID", mock.Anything, "media123").Return(media, nil)
			},
			expectedError: "forbidden",
		},
		{
			name:       "database error on get",
			uploaderID: "user123",
			mediaID:    "media123",
			mockSetup: func(mockRepo *MockMediaRepository, mockStorage *MockMediaStorage) {
				mockRepo.On("GetByID", mock.Anything, "media123").Return(nil, errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockMediaRepository)
			mockStorage := new(MockMediaStorage)
			tt.mockSetup(mockRepo, mockStorage)

			uc := usecase.NewMediaUseCase(
				mockRepo,
				mockStorage,
				&usecase.Config{},
				logger,
			)

			err := uc.Delete(context.Background(), tt.uploaderID, tt.mediaID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
		})
	}
}
