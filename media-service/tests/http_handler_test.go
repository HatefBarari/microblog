package tests

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HatefBarari/microblog-media/internal/domain"
	"github.com/HatefBarari/microblog-media/internal/presenter"
	"github.com/HatefBarari/microblog-media/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock use case
type MockMediaUseCase struct {
	mock.Mock
}

func (m *MockMediaUseCase) Upload(ctx context.Context, uploaderID string, file *multipart.FileHeader) (*usecase.UploadResponse, error) {
	args := m.Called(ctx, uploaderID, file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.UploadResponse), args.Error(1)
}

func (m *MockMediaUseCase) GetByID(ctx context.Context, id string) (*usecase.MediaResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.MediaResponse), args.Error(1)
}

func (m *MockMediaUseCase) List(ctx context.Context, filter usecase.ListFilter) (*usecase.ListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.ListResponse), args.Error(1)
}

func (m *MockMediaUseCase) Delete(ctx context.Context, uploaderID, id string) error {
	args := m.Called(ctx, uploaderID, id)
	return args.Error(0)
}

func setupTestHandler() (*presenter.HTTPHandler, *MockMediaUseCase) {
	mockUC := new(MockMediaUseCase)
	handler := presenter.NewHTTPHandler(mockUC)
	return handler, mockUC
}

func TestHTTPHandler_Upload(t *testing.T) {
	handler, mockUC := setupTestHandler()
	
	tests := []struct {
		name           string
		userID         string
		role           string
		file           *multipart.FileHeader
		mockSetup      func(*MockMediaUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful upload by admin",
			userID: "user123",
			role:   "admin",
			file: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024,
				Header:   map[string][]string{"Content-Type": {"image/jpeg"}},
			},
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("Upload", mock.Anything, "user123", mock.AnythingOfType("*multipart.FileHeader")).
					Return(&usecase.UploadResponse{
						Media: &usecase.MediaResponse{
							ID:           "media123",
							UploaderID:   "user123",
							Filename:     "test.jpg",
							OriginalName: "test.jpg",
							Type:         string(domain.MediaTypeImage),
						},
						URL: "http://localhost:8083/uploads/test.jpg",
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "successful upload by author",
			userID: "user123",
			role:   "author",
			file: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024,
				Header:   map[string][]string{"Content-Type": {"image/jpeg"}},
			},
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("Upload", mock.Anything, "user123", mock.AnythingOfType("*multipart.FileHeader")).
					Return(&usecase.UploadResponse{
						Media: &usecase.MediaResponse{
							ID:           "media123",
							UploaderID:   "user123",
							Filename:     "test.jpg",
							OriginalName: "test.jpg",
							Type:         string(domain.MediaTypeImage),
						},
						URL: "http://localhost:8083/uploads/test.jpg",
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "forbidden - user role",
			userID: "user123",
			role:   "user",
			file: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024,
				Header:   map[string][]string{"Content-Type": {"image/jpeg"}},
			},
			mockSetup:      func(mockUC *MockMediaUseCase) {},
			expectedStatus: http.StatusForbidden,
			expectedError:  "insufficient permissions",
		},
		{
			name:   "forbidden - guest role",
			userID: "user123",
			role:   "guest",
			file: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024,
				Header:   map[string][]string{"Content-Type": {"image/jpeg"}},
			},
			mockSetup:      func(mockUC *MockMediaUseCase) {},
			expectedStatus: http.StatusForbidden,
			expectedError:  "insufficient permissions",
		},
		{
			name:   "no file uploaded",
			userID: "user123",
			role:   "admin",
			file:   nil,
			mockSetup: func(mockUC *MockMediaUseCase) {
				// No repository calls expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "no file uploaded",
		},
		{
			name:   "use case error",
			userID: "user123",
			role:   "admin",
			file: &multipart.FileHeader{
				Filename: "test.jpg",
				Size:     1024,
				Header:   map[string][]string{"Content-Type": {"image/jpeg"}},
			},
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("Upload", mock.Anything, "user123", mock.AnythingOfType("*multipart.FileHeader")).
					Return(nil, errors.New("file too large"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "file too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUC.ExpectedCalls = nil
			tt.mockSetup(mockUC)

			// Setup request
			req := httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", nil)
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)
			c.Set("role", tt.role)

			// Add file to form if provided
			if tt.file != nil {
				// Create a simple form with file
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("file", tt.file.Filename)
				part.Write([]byte("fake image data"))
				writer.Close()
				
				req = httptest.NewRequest(http.MethodPost, "/api/v1/media/upload", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				rec = httptest.NewRecorder()
				
				c = e.NewContext(req, rec)
				c.Set("userID", tt.userID)
				c.Set("role", tt.role)
			}

			// Execute
			err := handler.Upload(c)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockUC.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_GetByID(t *testing.T) {
	handler, mockUC := setupTestHandler()
	
	tests := []struct {
		name           string
		mediaID        string
		mockSetup      func(*MockMediaUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful get by ID",
			mediaID: "media123",
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("GetByID", mock.Anything, "media123").
					Return(&usecase.MediaResponse{
						ID:           "media123",
						UploaderID:   "user123",
						Filename:     "test.jpg",
						OriginalName: "test.jpg",
						Type:         string(domain.MediaTypeImage),
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "media not found",
			mediaID: "nonexistent",
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("GetByID", mock.Anything, "nonexistent").
					Return(nil, errors.New("media not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "media not found",
		},
		{
			name:    "use case error",
			mediaID: "media123",
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("GetByID", mock.Anything, "media123").
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUC.ExpectedCalls = nil
			tt.mockSetup(mockUC)

			// Setup request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/media/"+tt.mediaID, nil)
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.mediaID)

			// Execute
			err := handler.GetByID(c)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockUC.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_List(t *testing.T) {
	handler, mockUC := setupTestHandler()
	
	tests := []struct {
		name           string
		userID         string
		queryParams    map[string]string
		mockSetup      func(*MockMediaUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful list",
			userID: "user123",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "10",
			},
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("List", mock.Anything, mock.MatchedBy(func(filter usecase.ListFilter) bool {
					return filter.UploaderID == "user123" &&
						filter.Page == 1 &&
						filter.PageSize == 10
				})).Return(&usecase.ListResponse{
					Media: []*usecase.MediaResponse{
						{
							ID:           "media1",
							UploaderID:   "user123",
							Filename:     "test1.jpg",
							OriginalName: "test1.jpg",
							Type:         string(domain.MediaTypeImage),
						},
						{
							ID:           "media2",
							UploaderID:   "user123",
							Filename:     "test2.jpg",
							OriginalName: "test2.jpg",
							Type:         string(domain.MediaTypeImage),
						},
					},
					Total: 2,
					Page:  1,
					Size:  2,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "empty list",
			userID: "user123",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "10",
			},
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("List", mock.Anything, mock.AnythingOfType("usecase.ListFilter")).
					Return(&usecase.ListResponse{
						Media: []*usecase.MediaResponse{},
						Total: 0,
						Page:  1,
						Size:  0,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "use case error",
			userID: "user123",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "10",
			},
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("List", mock.Anything, mock.AnythingOfType("usecase.ListFilter")).
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUC.ExpectedCalls = nil
			tt.mockSetup(mockUC)

			// Setup request with query parameters
			req := httptest.NewRequest(http.MethodGet, "/api/v1/media", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)

			// Execute
			err := handler.List(c)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockUC.AssertExpectations(t)
		})
	}
}

func TestHTTPHandler_Delete(t *testing.T) {
	handler, mockUC := setupTestHandler()
	
	tests := []struct {
		name           string
		userID         string
		mediaID        string
		mockSetup      func(*MockMediaUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful delete",
			userID:  "user123",
			mediaID: "media123",
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("Delete", mock.Anything, "user123", "media123").Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:    "media not found",
			userID:  "user123",
			mediaID: "nonexistent",
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("Delete", mock.Anything, "user123", "nonexistent").
					Return(errors.New("media not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "media not found",
		},
		{
			name:    "forbidden",
			userID:  "user456",
			mediaID: "media123",
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("Delete", mock.Anything, "user456", "media123").
					Return(errors.New("forbidden"))
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "forbidden",
		},
		{
			name:    "use case error",
			userID:  "user123",
			mediaID: "media123",
			mockSetup: func(mockUC *MockMediaUseCase) {
				mockUC.On("Delete", mock.Anything, "user123", "media123").
					Return(errors.New("database error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockUC.ExpectedCalls = nil
			tt.mockSetup(mockUC)

			// Setup request
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/media/"+tt.mediaID, nil)
			rec := httptest.NewRecorder()

			// Setup echo context
			e := echo.New()
			c := e.NewContext(req, rec)
			c.Set("userID", tt.userID)
			c.SetParamNames("id")
			c.SetParamValues(tt.mediaID)

			// Execute
			err := handler.Delete(c)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockUC.AssertExpectations(t)
		})
	}
}
