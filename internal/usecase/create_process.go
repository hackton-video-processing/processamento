package usecase

import (
	"context"
	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/mysql"
)

type CreateProcess struct {
	repository mysql.Repository
}

func NewCreateProcess(repository *mysql.Repository) *CreateProcess {
	return &CreateProcess{
		repository: *repository,
	}
}

func (c *CreateProcess) Execute(ctx context.Context, videoProcessing videoprocessing.VideoProcessing) error {
	err := c.repository.Create(ctx, videoProcessing)
	if err != nil {
		return err
	}

	return nil
}
