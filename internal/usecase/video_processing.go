package usecase

import (
	"fmt"
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
		notificationAPI string
	}

	VideoProcessingRequest struct {
		VideoKeys []string `json:"video_keys"`
	}
)

func NewVideoProcessing(s3Client *s3.S3Client, repository *mysql.Repository) *VideoProcessing {
	return &VideoProcessing{
		s3Client:   s3Client,
		repository: *repository,
	}
}

func (v *VideoProcessing) Execute(req VideoProcessingRequest) error {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 3) // Limita a 3 vídeos sendo processados simultaneamente

	for _, videoKey := range req.VideoKeys {
		wg.Add(1)
		go func(videoKey string) {
			defer wg.Done()
			sem <- struct{}{}        // Bloqueia se atingir o limite de concorrência
			defer func() { <-sem }() // Libera o slot após processamento

			// Baixa o vídeo do S3
			videoOutput, err := v.s3Client.GetVideo(videoKey)
			if err != nil {
				log.Printf("Erro ao baixar vídeo %s: %v", videoKey, err)
				return
			}

			// Salva o vídeo localmente
			localVideoPath := fmt.Sprintf("/tmp/%s", filepath.Base(videoKey))
			file, err := os.Create(localVideoPath)
			if err != nil {
				log.Printf("Erro ao criar arquivo local para %s: %v", videoKey, err)
				return
			}
			defer file.Close()
			io.Copy(file, videoOutput.Body)

			// Processa o vídeo (extrai frames)
			outputDir := fmt.Sprintf("/tmp/frames_%s", filepath.Base(videoKey))
			os.MkdirAll(outputDir, os.ModePerm)
			if err := processVideo(localVideoPath, outputDir); err != nil {
				log.Printf("Erro ao processar vídeo %s: %v", videoKey, err)
				return
			}

			// Compacta os frames
			zipPath := fmt.Sprintf("%s.zip", outputDir)
			if err := zip.ZipDirectory(outputDir, zipPath); err != nil {
				log.Printf("Erro ao compactar imagens de %s: %v", videoKey, err)
				return
			}

			// Faz upload do ZIP para o S3
			zipKey := fmt.Sprintf("processed/%s.zip", filepath.Base(videoKey))
			if err := v.s3Client.UploadZippedVideo(zipKey, zipPath); err != nil {
				log.Printf("Erro ao fazer upload do zip %s: %v", zipKey, err)
				return
			}

			log.Printf("Processamento de %s concluído com sucesso", videoKey)
		}(videoKey)
	}

	wg.Wait() // Aguarda todos os vídeos serem processados
	log.Println("Todos os vídeos foram processados com sucesso!")
	return nil
}

func processVideo(videoPath string, outputDir string) error {
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=1", filepath.Join(outputDir, "frame-%03d.jpg"))
	return cmd.Run()
}
