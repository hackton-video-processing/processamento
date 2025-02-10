package kafka

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hackton-video-processing/processamento/internal/usecase"
)

type (
	videoProcessingUseCase interface {
		Execute(context context.Context, request usecase.VideoProcessingRequest) error
	}

	VideoProcessingHandler struct {
		videoProcessingUseCase videoProcessingUseCase
	}
)

func NewVideoProcessingHandler(videoProcessingUseCase videoProcessingUseCase) *VideoProcessingHandler {
	return &VideoProcessingHandler{
		videoProcessingUseCase: videoProcessingUseCase,
	}
}

func (v *VideoProcessingHandler) VideoProcessing(w http.ResponseWriter, r *http.Request) {
	var req usecase.VideoProcessingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := v.videoProcessingUseCase.Execute(r.Context(), req); err != nil {
		http.Error(w, "Error processing video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Video processing completed"))
}
