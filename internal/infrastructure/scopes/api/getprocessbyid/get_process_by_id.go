package getprocessbyid

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
	"net/http"
)

type (
	getProcessByIDUseCase interface {
		Execute(ctx context.Context, videoProcessingID string) (videoprocessing.VideoProcessing, error)
	}

	GetProcessByIDHandler struct {
		getProcessByIDUseCase getProcessByIDUseCase
	}
)

func NewGetProcessByIDHandler(getProcessByIDUseCase getProcessByIDUseCase) *GetProcessByIDHandler {
	return &GetProcessByIDHandler{
		getProcessByIDUseCase: getProcessByIDUseCase,
	}
}

func (g *GetProcessByIDHandler) GetProcessByID(w http.ResponseWriter, r *http.Request) {
	processID := chi.URLParam(r, "id")
	if processID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	process, err := g.getProcessByIDUseCase.Execute(ctx, processID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(process)
}
