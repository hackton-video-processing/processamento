package s3

import (
	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func BootstrapS3(appConfig config.AppConfig) *S3Client {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(appConfig.S3Config.Region)}))
	s3Client := s3.New(sess)

	return NewS3Client(appConfig, s3Client)
}
