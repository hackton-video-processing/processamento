package createprocess

import (
	"context"
	"encoding/json"
	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
	"net/http"
)

type (
	createProcessUseCase interface {
		Execute(ctx context.Context, videoProcessing videoprocessing.VideoProcessing) error
	}

	CreateProcessHandler struct {
		createProcessUseCase createProcessUseCase
	}

	request struct {
		Files []string `json:"files"`
	}
)

func NewCreateProcessHandler(createProcessUseCase createProcessUseCase) *CreateProcessHandler {
	return &CreateProcessHandler{
		createProcessUseCase: createProcessUseCase,
	}
}

func (h *CreateProcessHandler) CreateProcess(w http.ResponseWriter, r *http.Request) {
	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Files) == 0 {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.createProcessUseCase.Execute(ctx, videoprocessing.NewVideoProcessing(req.Files)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "process created successfully"})
}
