package catalog

import (
	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/mysql"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/s3"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"
	"github.com/hackton-video-processing/processamento/internal/usecase"
	"github.com/hackton-video-processing/processamento/pkg/notificationapi"
)

type UseCase struct {
	appConfig config.AppConfig
	s3Client  *s3.S3Client
	*mysql.Repository
}

func (u UseCase) Health() (*usecase.HealthCheck, error) {
	return usecase.NewHealthCheck(), nil
}

func (u UseCase) CreateProcess() (*usecase.CreateProcess, error) {
	return usecase.NewCreateProcess(u.Repository), nil
}

func (u UseCase) GetProcessByID() (*usecase.GetProcessByID, error) {
	return usecase.NewGetProcessByID(u.Repository), nil
}

func (u UseCase) Process() (*usecase.VideoProcessing, error) {
	notificationAPI, err := notificationapi.NewNotificationService(u.appConfig.NotificationAPIConfig.BaseURL, u.appConfig.NotificationAPIConfig.Endpoint)
	if err != nil {
		return nil, err
	}

	return usecase.NewVideoProcessing(u.s3Client, u.Repository, u.appConfig, notificationAPI), nil
}

func New(appConfig config.AppConfig, s3Client *s3.S3Client, repository *mysql.Repository) UseCase {
	return UseCase{
		appConfig:  appConfig,
		s3Client:   s3Client,
		Repository: repository,
	}
}
