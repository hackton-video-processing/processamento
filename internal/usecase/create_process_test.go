package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
	"github.com/hackton-video-processing/processamento/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository é um mock do mysql.Repository
type MockCreateRepository struct {
	mock.Mock
}

func (m *MockCreateRepository) Create(ctx context.Context, videoProcessing videoprocessing.VideoProcessing) (string, error) {
	args := m.Called(ctx, videoProcessing)
	return args.String(0), args.Error(1)
}

func TestCreateProcess_Execute_Success(t *testing.T) {
	mockRepo := new(MockCreateRepository)
	createProcess := usecase.NewCreateProcess(mockRepo)

	video := videoprocessing.VideoProcessing{
		ID:     "process123",
		Status: videoprocessing.Created,
		Files: []videoprocessing.File{
			{ID: "file1", Name: "video.mp4"},
		},
	}

	mockRepo.On("Create", mock.Anything, video).Return("process123", nil)

	processID, err := createProcess.Execute(context.Background(), video)

	assert.NoError(t, err)
	assert.Equal(t, "process123", processID)
	mockRepo.AssertCalled(t, "Create", mock.Anything, video)
}

func TestCreateProcess_Execute_Error(t *testing.T) {
	mockRepo := new(MockCreateRepository)
	createProcess := usecase.NewCreateProcess(mockRepo)

	video := videoprocessing.VideoProcessing{
		ID:     "process123",
		Status: videoprocessing.Created,
		Files: []videoprocessing.File{
			{ID: "file1", Name: "video.mp4"},
		},
	}

	mockRepo.On("Create", mock.Anything, video).Return("", errors.New("database error"))

	processID, err := createProcess.Execute(context.Background(), video)

	assert.Error(t, err)
	assert.Equal(t, "", processID)
	assert.EqualError(t, err, "database error")
	mockRepo.AssertCalled(t, "Create", mock.Anything, video)
}
