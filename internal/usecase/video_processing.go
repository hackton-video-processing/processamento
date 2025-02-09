package usecase

import (
	"context"
	"fmt"
	"github.com/hackton-video-processing/processamento/internal/domain/videoprocessing"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/mysql"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/s3"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"
	"github.com/hackton-video-processing/processamento/pkg/notificationapi"
	"github.com/hackton-video-processing/processamento/pkg/zip"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type (
	VideoProcessing struct {
		s3Client        *s3.S3Client
		repository      mysql.Repository
		appConfig       config.AppConfig
		notificationAPI *notificationapi.NotificationService
	}

	VideoProcessingRequest struct {
		Email       string `json:"email"`
		ProcessedID string `json:"processed_id"`
	}
)

func NewVideoProcessing(s3Client *s3.S3Client, repository *mysql.Repository, appConfig config.AppConfig, notificationAPI *notificationapi.NotificationService) *VideoProcessing {
	return &VideoProcessing{
		s3Client:        s3Client,
		repository:      *repository,
		appConfig:       appConfig,
		notificationAPI: notificationAPI,
	}
}

func (v *VideoProcessing) Execute(ctx context.Context, req VideoProcessingRequest) error {
	var (
		wg  sync.WaitGroup
		sem = make(chan struct{}, v.appConfig.VideoProcessingConfig.MaxVideoProcessing)
	)

	videoProcessing, err := v.repository.GetProcessByID(ctx, req.ProcessedID)
	if err != nil {
		return err
	}

	err = v.repository.UpdateStatusByID(ctx, req.ProcessedID, "in_progress")
	if err != nil {
		log.Fatalf("failed to update status: %v", err)
	}

	fileProcessed := len(videoProcessing.Files)
	for _, s3file := range videoProcessing.Files {
		wg.Add(1)
		go func(s3file videoprocessing.File) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			uploadedPath := fmt.Sprintf("%s%s", v.appConfig.S3Config.UploadPath, s3file.Name)
			videoOutput, goErr := v.s3Client.GetVideo(uploadedPath)
			if goErr != nil {
				log.Printf("Erro ao baixar vídeo %s: %v", s3file.Name, goErr)
				err = goErr
				return
			}

			localVideoPath := fmt.Sprintf("/videos/%s", filepath.Base(s3file.Name))
			file, goErr := os.Create(localVideoPath)
			if goErr != nil {
				log.Printf("Erro ao criar arquivo local para %s: %v", file.Name, goErr)
				err = goErr
				return
			}
			defer file.Close()
			io.Copy(file, videoOutput.Body)

			outputDir := fmt.Sprintf("/tmp/frames_%s", filepath.Base(s3file.Name))
			os.MkdirAll(outputDir, os.ModePerm)
			if goErr = processVideo(localVideoPath, outputDir); goErr != nil {
				log.Printf("Erro ao processar vídeo %s: %v", s3file.Name, goErr)
				err = goErr
				return
			}

			zipPath := fmt.Sprintf("%s.zip", outputDir)
			if goErr = zip.ZipDirectory(outputDir, zipPath); goErr != nil {
				log.Printf("Erro ao compactar imagens de %s: %v", s3file.Name, goErr)
				err = goErr
				return
			}

			destinationPath := fmt.Sprintf("%s%s.zip", v.appConfig.S3Config.DownloadPath, filepath.Base(s3file.Name))
			if goErr = v.s3Client.UploadZippedVideo(destinationPath, zipPath); goErr != nil {
				log.Printf("Erro ao fazer upload do zip %s: %v", destinationPath, goErr)
				err = goErr
				return
			}

			if goErr = v.s3Client.DeleteVideo(s3file.Name); goErr != nil {
				log.Printf("Erro ao excluir vídeo original %s: %v", s3file.Name, goErr)
				err = goErr
				return
			}

			log.Printf("Processamento de %s concluído com sucesso e vídeo original excluído", s3file.Name)

			fileProcessed++
		}(s3file)
	}

	wg.Wait()

	if err != nil {
		log.Println("Nem todos os vídeos foram processados com sucesso.")
		err = v.repository.UpdateStatusByID(ctx, req.ProcessedID, "failed")
		if err != nil {
			return fmt.Errorf("failed to update status to failed: %w", err)
		}

		err = v.notificationAPI.SendNotification(req.Email, "Tivemos problema ao processar sua requisição")
		if err != nil {
			return fmt.Errorf("failed to send notification: %w", err)
		}

		return err
	}

	if fileProcessed == len(videoProcessing.Files) {
		log.Println("Todos os vídeos foram processados com sucesso!")
		err = v.repository.UpdateStatusByID(ctx, req.ProcessedID, "completed")
		if err != nil {
			return fmt.Errorf("failed to update status to completed: %w", err)
		}

		err := v.notificationAPI.SendNotification(req.Email, "O vídeo já está disponível para download.")
		if err != nil {
			return fmt.Errorf("failed to send notification: %w", err)
		}
	}

	return nil
}

func processVideo(videoPath string, outputDir string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=1", filepath.Join(outputDir, "frame-%03d.jpg"))
	return cmd.Run()
}
