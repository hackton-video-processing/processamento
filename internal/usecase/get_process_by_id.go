package usecase

import (
	"context"
	"errors"
	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
)

type GetProcessByID struct {
	repository GetProcessRepository
}

type GetProcessRepository interface {
	GetProcessByID(ctx context.Context, processID string) (videoprocessing.VideoProcessing, error)
}

func NewGetProcessByID(repository GetProcessRepository) *GetProcessByID {
	return &GetProcessByID{
		repository: repository,
	}
}

func (g *GetProcessByID) Execute(ctx context.Context, videoProcessingID string) (videoprocessing.VideoProcessing, error) {
	videoProcess, err := g.repository.GetProcessByID(ctx, videoProcessingID)
	if err != nil {
		if errors.Is(err, videoprocessing.ErrVideoProcessingNotFound) {
			return videoprocessing.VideoProcessing{}, nil
		}
		return videoprocessing.VideoProcessing{}, err
	}

	return videoProcess, nil
}
