package healthcheck_test

import (
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/api/healthcheck"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHealthCheckUseCase struct {
	checkFunc func() string
}

func (m *mockHealthCheckUseCase) Check() string {
	return m.checkFunc()
}

func TestHealthCheckHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Healthy response",
			mockResponse:   "OK",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Custom health response",
			mockResponse:   "Service is running",
			expectedStatus: http.StatusOK,
			expectedBody:   "Service is running",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &mockHealthCheckUseCase{
				checkFunc: func() string {
					return tt.mockResponse
				},
			}

			handler := healthcheck.NewHealthCheckHandler(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			handler.HealthCheck(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedBody, rec.Body.String())
		})
	}
}
