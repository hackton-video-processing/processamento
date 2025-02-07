package mysql

import (
	"context"
	"errors"
	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"

	"gorm.io/gorm"
)

type (
	Repository struct {
		db *gorm.DB
	}
)

func NewMySQLRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) GetProcessByID(ctx context.Context, processID string) (videoprocessing.VideoProcessing, error) {
	var videoProcessing ProcessMySQL
	result := r.db.WithContext(ctx).
		Preload("Files").
		Where("process_id = ?", processID).
		First(&videoProcessing)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return videoprocessing.VideoProcessing{}, videoprocessing.ErrVideoProcessingNotFound
		}

		return videoprocessing.VideoProcessing{}, result.Error
	}

	return toDomain(videoProcessing), nil
}

func (r Repository) Create(ctx context.Context, videoProcessing videoprocessing.VideoProcessing) (string, error) {
	process := fromDomain(videoProcessing)

	result := r.db.WithContext(ctx).Create(&process)
	if result.Error != nil {
		return "", result.Error
	}

	return process.ID, nil
}

func (r Repository) UpdateStatusByID(ctx context.Context, processID string, status string) error {
	result := r.db.WithContext(ctx).
		Model(&ProcessMySQL{}).
		Where("process_id = ?", processID).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return videoprocessing.ErrVideoProcessingNotFound
	}

	return nil
}
