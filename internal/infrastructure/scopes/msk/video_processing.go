package msk

import (
	"encoding/json"
	"net/http"

	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"
	"github.com/hackton-video-processing/processamento/internal/usecase"
)

type (
	videoProcessingUseCase interface {
		Execute(request usecase.VideoProcessingRequest) error
	}

	VideoProcessingHandler struct {
		appConfig              config.AppConfig
		videoProcessingUseCase videoProcessingUseCase
	}
)

func NewVideoProcessingHandler(appConfig config.AppConfig, videoProcessingUseCase videoProcessingUseCase) *VideoProcessingHandler {
	return &VideoProcessingHandler{
		appConfig:              appConfig,
		videoProcessingUseCase: videoProcessingUseCase,
	}
}

// VideoProcessing processa os vídeos
func (v *VideoProcessingHandler) VideoProcessing(w http.ResponseWriter, r *http.Request) {
	var req usecase.VideoProcessingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := v.videoProcessingUseCase.Execute(req); err != nil {
		http.Error(w, "Error processing video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Video processing completed"))
}
