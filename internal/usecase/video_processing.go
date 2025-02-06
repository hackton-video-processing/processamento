package usecase

import (
	"context"
	"fmt"
	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/mysql"
	"github.com/hackton-video-processing/processamento/pkg/zip"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/s3"
)

type (
	VideoProcessing struct {
		s3Client        *s3.S3Client
		repository      mysql.Repository
		maxConcurrency  int
		notificationAPI string
	}

	VideoProcessingRequest struct {
		Email       string `json:"email"`
		ProcessedID string `json:"processed_id"`
	}
)

func NewVideoProcessing(s3Client *s3.S3Client, repository *mysql.Repository, maxConcurrency int) *VideoProcessing {
	return &VideoProcessing{
		s3Client:       s3Client,
		repository:     *repository,
		maxConcurrency: maxConcurrency,
	}
}

func (v *VideoProcessing) Execute(ctx context.Context, req VideoProcessingRequest) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, v.maxConcurrency)

	videoProcessing, err := v.repository.GetProcessByID(ctx, req.ProcessedID)
	if err != nil {
		return err
	}

	for _, videoKey := range videoProcessing.Files {
		wg.Add(1)
		go func(videoKey videoprocessing.File) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			videoOutput, err := v.s3Client.GetVideo(videoKey.Name)
			if err != nil {
				log.Printf("Erro ao baixar vídeo %s: %v", videoKey, err)
				return
			}

			localVideoPath := fmt.Sprintf("/tmp/%s", filepath.Base(videoKey.Name))
			file, err := os.Create(localVideoPath)
			if err != nil {
				log.Printf("Erro ao criar arquivo local para %s: %v", videoKey, err)
				return
			}
			defer file.Close()
			io.Copy(file, videoOutput.Body)

			outputDir := fmt.Sprintf("/tmp/frames_%s", filepath.Base(videoKey.Name))
			os.MkdirAll(outputDir, os.ModePerm)
			if err := processVideo(localVideoPath, outputDir); err != nil {
				log.Printf("Erro ao processar vídeo %s: %v", videoKey, err)
				return
			}

			zipPath := fmt.Sprintf("%s.zip", outputDir)
			if err := zip.ZipDirectory(outputDir, zipPath); err != nil {
				log.Printf("Erro ao compactar imagens de %s: %v", videoKey, err)
				return
			}

			zipKey := fmt.Sprintf("processed/%s.zip", filepath.Base(videoKey.Name))
			if err := v.s3Client.UploadZippedVideo(zipKey, zipPath); err != nil {
				log.Printf("Erro ao fazer upload do zip %s: %v", zipKey, err)
				return
			}

			log.Printf("Processamento de %s concluído com sucesso", videoKey)
		}(videoKey)
	}

	wg.Wait()
	log.Println("Todos os vídeos foram processados com sucesso!")
	return nil
}

func processVideo(videoPath string, outputDir string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=1", filepath.Join(outputDir, "frame-%03d.jpg"))
	return cmd.Run()
}
