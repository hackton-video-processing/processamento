package kafka_test

import (
	"context"
	"errors"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/kafka"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hackton-video-processing/processamento/internal/usecase"
	"github.com/stretchr/testify/assert"
)

type mockVideoProcessingUseCase struct {
	executeFunc func(ctx context.Context, req usecase.VideoProcessingRequest) error
}

func (m *mockVideoProcessingUseCase) Execute(ctx context.Context, req usecase.VideoProcessingRequest) error {
	return m.executeFunc(ctx, req)
}

func TestVideoProcessingHandler_VideoProcessing(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid request payload",
			requestBody:    "{invalid-json}",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request payload\n",
		},
		{
			name:           "Error during video processing",
			requestBody:    `{"video_id":"1234", "source_url":"http://example.com/video.mp4"}`,
			mockError:      errors.New("processing error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Error processing video: processing error\n",
		},
		{
			name:           "Successful video processing",
			requestBody:    `{"video_id":"1234", "source_url":"http://example.com/video.mp4"}`,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Video processing completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &mockVideoProcessingUseCase{
				executeFunc: func(ctx context.Context, req usecase.VideoProcessingRequest) error {
					return tt.mockError
				},
			}

			handler := kafka.NewVideoProcessingHandler(mockUseCase)

			req := httptest.NewRequest(http.MethodPost, "/process", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.VideoProcessing(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedBody, rec.Body.String())
		})
	}
}
