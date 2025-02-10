package usecase

import (
	"context"
	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
)

type CreateProcess struct {
	repository CreateProcessRepository
}

type CreateProcessRepository interface {
	Create(ctx context.Context, videoProcessing videoprocessing.VideoProcessing) (string, error)
}

func NewCreateProcess(repository CreateProcessRepository) *CreateProcess {
	return &CreateProcess{
		repository: repository,
	}
}

func (c *CreateProcess) Execute(ctx context.Context, videoProcessing videoprocessing.VideoProcessing) (string, error) {
	processID, err := c.repository.Create(ctx, videoProcessing)
	if err != nil {
		return "", err
	}

	return processID, nil
}
