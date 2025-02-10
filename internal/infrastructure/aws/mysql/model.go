package mysql

import (
	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
	"time"
)

type ProcessMySQL struct {
	ID        string    `gorm:"column:process_id;primaryKey;not null"`
	Files     []File    `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Status    string    `gorm:"column:status;not null"`
	CreatedAt time.Time `gorm:"column:create_at;autoCreateTime"`
}

func (ProcessMySQL) TableName() string {
	return "process"
}

type File struct {
	ID        string `gorm:"column:file_id;primaryKey;not null"`
	Name      string `gorm:"column:file_name;not null"`
	Link      string `gorm:"column:link;"`
	ProcessID string `gorm:"column:process_id;index;not null"`
}

func (File) TableName() string {
	return "file"
}

func fromDomain(videoProcess videoprocessing.VideoProcessing) *ProcessMySQL {
	model := &ProcessMySQL{
		ID:     videoProcess.ID,
		Status: string(videoProcess.Status),
	}

	for _, file := range videoProcess.Files {
		f := &File{
			ID:        file.ID,
			Name:      file.Name,
			Link:      file.Link,
			ProcessID: videoProcess.ID,
		}

		model.Files = append(model.Files, *f)
	}

	return model
}

func toDomain(sql ProcessMySQL) videoprocessing.VideoProcessing {
	process := videoprocessing.VideoProcessing{
		ID:        sql.ID,
		Status:    videoprocessing.Status(sql.Status),
		CreatedAt: sql.CreatedAt,
	}

	for _, file := range sql.Files {
		f := videoprocessing.File{
			ID:   file.ID,
			Name: file.Name,
			Link: file.Link,
		}

		process.Files = append(process.Files, f)
	}

	return process
}
