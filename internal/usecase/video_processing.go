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
		ProcessedID string `json:"processId"`
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

	err = v.repository.UpdateStatusByProcessID(ctx, videoProcessing.ID, string(videoprocessing.Processing))
	if err != nil {
		log.Fatalf("error updating status: %v", err)
	}

	fileProcessed := len(videoProcessing.Files)
	for _, videoProcessingFile := range videoProcessing.Files {
		wg.Add(1)
		func(videoProcessingFile videoprocessing.File) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			s3uploadedPath := fmt.Sprintf("%s%s", v.appConfig.S3Config.UploadPath, videoProcessingFile.Name)
			s3File, goErr := v.s3Client.GetVideo(s3uploadedPath)
			if goErr != nil {
				log.Printf("error downloading S3 localFile %s: %v", videoProcessingFile.Name, goErr)
				err = goErr
				return
			}

			localDir := "./videos"
			if goErr = os.MkdirAll(localDir, os.ModePerm); err != nil {
				log.Printf("error creating local dir %s: %v", localDir, goErr)
				return
			}

			localVideoPath := fmt.Sprintf("%s/%s", localDir, filepath.Base(videoProcessingFile.Name))
			localFile, goErr := os.Create(localVideoPath)
			if goErr != nil {
				log.Printf("error creating local localFile for %s: %v", localVideoPath, goErr)
				err = goErr
				return
			}
			defer localFile.Close()

			io.Copy(localFile, s3File.Body)

			tempOutputDir := fmt.Sprintf("./tmp/frames_%s", filepath.Base(videoProcessingFile.Name))
			if err := os.MkdirAll(tempOutputDir, os.ModePerm); err != nil {
				log.Printf("error creating temp output dir %s: %v", videoProcessingFile.Name, err)
				return
			}

			if goErr = processVideo(localVideoPath, tempOutputDir); goErr != nil {
				log.Printf("error processing video %s: %v", videoProcessingFile.Name, goErr)
				err = goErr
				return
			}

			zipPath := fmt.Sprintf("%s.zip", tempOutputDir)
			if goErr = zip.ZipDirectory(tempOutputDir, zipPath); goErr != nil {
				log.Printf("error zipping video %s: %v", videoProcessingFile.Name, goErr)
				err = goErr
				return
			}

			s3DestinationPath := fmt.Sprintf("%s%s.zip", v.appConfig.S3Config.DownloadPath, filepath.Base(videoProcessingFile.Name))
			if goErr = v.s3Client.UploadZippedVideo(s3DestinationPath, zipPath); goErr != nil {
				log.Printf("error uploading zip file %s: %v", s3DestinationPath, goErr)
				err = goErr
				return
			}

			goErr = v.repository.UpdateFileByID(ctx, videoProcessingFile.ID, fmt.Sprintf(
				"https://%s.s3.%s.amazonaws.com/%s%s.zip",
				v.appConfig.S3Config.S3Bucket,
				v.appConfig.S3Config.Region,
				v.appConfig.S3Config.DownloadPath,
				videoProcessingFile.Name))
			if goErr != nil {
				log.Printf("error updating status: %v", goErr)
				err = goErr
				return
			}

			if goErr = v.s3Client.DeleteVideo(videoProcessingFile.Name); goErr != nil {
				log.Printf("error deleting original s3 video %s: %v", videoProcessingFile.Name, goErr)
				err = goErr
				return
			}

			fileProcessed++
		}(videoProcessingFile)
	}

	wg.Wait()

	if err != nil {
		log.Println("one or more errors occurred.")

		err2 := v.repository.UpdateStatusByProcessID(ctx, req.ProcessedID, string(videoprocessing.Failed))
		if err2 != nil {
			return fmt.Errorf("error updating status: %w", err2)
		}

		err2 = v.notificationAPI.SendNotification(req.Email, "Tivemos problema ao processar sua requisição")
		if err2 != nil {
			return fmt.Errorf("error sending notification: %w", err2)
		}

		return err
	}

	if fileProcessed == len(videoProcessing.Files) {
		err = v.repository.UpdateStatusByProcessID(ctx, req.ProcessedID, string(videoprocessing.Processed))
		if err != nil {
			return fmt.Errorf("error updating status: %w", err)
		}

		err = v.notificationAPI.SendNotification(req.Email, "O vídeo já está disponível para download.")
		if err != nil {
			return fmt.Errorf("error sending notification: %w", err)
		}
	}

	log.Println("processing has been done successful")
	return nil
}

func processVideo(videoPath string, outputDir string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=1", filepath.Join(outputDir, "frame-%03d.jpg"))
	return cmd.Run()
}
