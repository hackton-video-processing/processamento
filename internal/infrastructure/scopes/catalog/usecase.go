package catalog

import (
	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/mysql"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/aws/s3"
	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"
	"github.com/hackton-video-processing/processamento/internal/usecase"
)

type UseCase struct {
	appConfig config.AppConfig
}

func (u UseCase) Health() (*usecase.HealthCheck, error) {
	return usecase.NewHealthCheck(), nil
}

func (u UseCase) CreateProcess() (*usecase.CreateProcess, error) {
	repository, err := u.BootstrapMySQLRepository()
	if err != nil {
		return nil, err
	}

	return usecase.NewCreateProcess(repository), nil
}

func (u UseCase) GetProcessByID() (*usecase.GetProcessByID, error) {
	repository, err := u.BootstrapMySQLRepository()
	if err != nil {
		return nil, err
	}

	return usecase.NewGetProcessByID(repository), nil
}

func (u UseCase) Process() (*usecase.VideoProcessing, error) {
	s3Client := u.BootstrapS3Client()

	repository, err := u.BootstrapMySQLRepository()
	if err != nil {
		return nil, err
	}

	return usecase.NewVideoProcessing(s3Client, repository), nil
}

func (u UseCase) BootstrapS3Client() *s3.S3Client {
	return s3.BootstrapS3(u.appConfig)
}

func (u UseCase) BootstrapMySQLRepository() (*mysql.Repository, error) {
	return mysql.BootstrapMySQLRepository(u.appConfig)
}

func New(appConfig config.AppConfig) UseCase {
	return UseCase{
		appConfig: appConfig,
	}
}
