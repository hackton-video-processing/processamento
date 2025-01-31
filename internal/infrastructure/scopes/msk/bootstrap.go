package msk

import (
	"fmt"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/scopes/catalog"
	"github.com/hackton-video-processing/processamento/pkg/once"
)

func BootStrapVideoProcessing(appConfig config.AppConfig) (*VideoProcessingHandler, error) {
	useCaseCatalog := catalog.New(appConfig)

	videoProcessingUecase, err := once.Call(useCaseCatalog.Process)
	if err != nil {
		return nil, fmt.Errorf("creating capture transaction use case: %w", err)
	}

	return NewVideoProcessingHandler(appConfig, videoProcessingUecase), nil
}
