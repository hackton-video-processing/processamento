package usecase_test

import (
	"context"
	"errors"
	"github.com/hackton-video-processing/processamento/internal/usecase"
	"testing"

	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGetProcessRepository struct {
	mock.Mock
}

func (m *MockGetProcessRepository) GetProcessByID(ctx context.Context, processID string) (videoprocessing.VideoProcessing, error) {
	args := m.Called(ctx, processID)
	return args.Get(0).(videoprocessing.VideoProcessing), args.Error(1)
}

func TestGetProcessByID_Execute_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockGetProcessRepository)
	expectedProcess := videoprocessing.VideoProcessing{ID: "123", Status: "completed"}

	// Configura o mock para retornar um processo de vídeo válido
	mockRepo.On("GetProcessByID", ctx, "123").Return(expectedProcess, nil)

	usecase := usecase.NewGetProcessByID(mockRepo)
	result, err := usecase.Execute(ctx, "123")

	assert.NoError(t, err)
	assert.Equal(t, expectedProcess, result)
	mockRepo.AssertExpectations(t)
}

func TestGetProcessByID_Execute_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockGetProcessRepository)

	// Configura o mock para retornar o erro de "processo não encontrado"
	mockRepo.On("GetProcessByID", ctx, "123").Return(videoprocessing.VideoProcessing{}, videoprocessing.ErrVideoProcessingNotFound)

	usecase := usecase.NewGetProcessByID(mockRepo)
	result, err := usecase.Execute(ctx, "123")

	assert.NoError(t, err)
	assert.Equal(t, videoprocessing.VideoProcessing{}, result)
	mockRepo.AssertExpectations(t)
}

func TestGetProcessByID_Execute_UnexpectedError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockGetProcessRepository)
	unexpectedError := errors.New("unexpected error")

	mockRepo.On("GetProcessByID", ctx, "123").Return(videoprocessing.VideoProcessing{}, unexpectedError)

	usecase := usecase.NewGetProcessByID(mockRepo)
	result, err := usecase.Execute(ctx, "123")

	assert.Error(t, err)
	assert.Equal(t, unexpectedError, err)
	assert.Equal(t, videoprocessing.VideoProcessing{}, result)
	mockRepo.AssertExpectations(t)
}
