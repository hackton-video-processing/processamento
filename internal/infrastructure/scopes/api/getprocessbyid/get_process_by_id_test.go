package getprocessbyid_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/api/getprocessbyid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockGetProcessByIDUseCase struct {
	executeFunc func(ctx context.Context, videoProcessingID string) (videoprocessing.VideoProcessing, error)
}

func (m *mockGetProcessByIDUseCase) Execute(ctx context.Context, videoProcessingID string) (videoprocessing.VideoProcessing, error) {
	return m.executeFunc(ctx, videoProcessingID)
}

func TestGetProcessByIDHandler_GetProcessByID(t *testing.T) {
	tests := []struct {
		name           string
		urlParam       string
		mockProcess    videoprocessing.VideoProcessing
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Error executing use case",
			urlParam:       "12345",
			mockError:      errors.New("error fetching process"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "error fetching process\n",
		},
		{
			name:     "Process found successfully",
			urlParam: "12345",
			mockProcess: videoprocessing.VideoProcessing{
				ID:     "12345",
				Status: "completed",
				Files: []videoprocessing.File{
					{ID: "file1", Name: "video.mp4", Link: "http://example.com/video.mp4"},
					{ID: "file2", Name: "thumbnail.jpg", Link: "http://example.com/thumbnail.jpg"},
				},
				CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"id": "12345",
				"status": "completed",
				"files": [
					{"id": "file1", "name": "video.mp4", "link": "http://example.com/video.mp4"},
					{"id": "file2", "name": "thumbnail.jpg", "link": "http://example.com/thumbnail.jpg"}
				],
				"created_at": "2025-01-15T12:00:00Z"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &mockGetProcessByIDUseCase{
				executeFunc: func(ctx context.Context, videoProcessingID string) (videoprocessing.VideoProcessing, error) {
					return tt.mockProcess, tt.mockError
				},
			}

			handler := getprocessbyid.NewGetProcessByIDHandler(mockUseCase)

			r := chi.NewRouter()
			r.Get("/process/{id}", handler.GetProcessByID)

			req := httptest.NewRequest(http.MethodGet, "/process/"+tt.urlParam, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				var expected, actual map[string]interface{}
				json.Unmarshal([]byte(tt.expectedBody), &expected)
				json.Unmarshal(rec.Body.Bytes(), &actual)
				assert.Equal(t, expected, actual)
			} else {
				assert.Equal(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}
