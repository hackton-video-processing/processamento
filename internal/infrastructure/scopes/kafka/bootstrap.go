package kafka

import (
	"fmt"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/catalog"
	"github.com/hackton-video-processing/processamento/pkg/once"
)

func BootStrapVideoProcessing(useCaseCatalog catalog.UseCase) (*VideoProcessingHandler, error) {
	videoProcessingUecase, err := once.Call(useCaseCatalog.Process)
	if err != nil {
		return nil, fmt.Errorf("creating capture transaction use case: %w", err)
	}

	return NewVideoProcessingHandler(videoProcessingUecase), nil
}
