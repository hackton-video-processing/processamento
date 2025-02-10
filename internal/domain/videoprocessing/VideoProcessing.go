package videoprocessing

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Status string

var (
	Created    Status = "created"
	Processing Status = "processing"
	Processed  Status = "processed"
	Failed     Status = "failed"

	ErrVideoProcessingNotFound = errors.New("video processing not found")
)

type VideoProcessing struct {
	ID        string    `json:"id"`
	Files     []File    `json:"files"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func NewVideoProcessing(files []string) VideoProcessing {
	videoProcessing := VideoProcessing{
		ID:     uuid.NewString(),
		Status: Created,
	}

	for _, fileName := range files {
		f := File{
			ID:   uuid.NewString(),
			Name: fileName,
		}
		videoProcessing.Files = append(videoProcessing.Files, f)
	}

	return videoProcessing
}
