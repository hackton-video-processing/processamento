package usecase_test

import (
	"github.com/hackton-video-processing/processamento/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Check(t *testing.T) {
	healthCheck := usecase.NewHealthCheck()
	result := healthCheck.Check()

	assert.Equal(t, "Ok", result, "HealthCheck should return 'Ok'")
}
