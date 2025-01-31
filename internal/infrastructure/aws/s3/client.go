package s3

import (
	"os"

	"github.com/hackton-video-processing/processamento/internal/infrastructure/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	s3Config config.S3Config
	client   *s3.S3
}

func NewS3Client(config config.AppConfig, client *s3.S3) *S3Client {
	return &S3Client{
		s3Config: config.S3Config,
		client:   client,
	}
}

func (s *S3Client) GetVideo(videoKey string) (*s3.GetObjectOutput, error) {
	video, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.s3Config.S3Bucket),
		Key:    aws.String(videoKey),
	})
	if err != nil {
		return nil, err
	}

	return video, nil
}

func (s *S3Client) UploadZippedVideo(zipKey, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.s3Config.S3Bucket),
		Key:    aws.String(zipKey),
		Body:   file,
	})

	return err
}
