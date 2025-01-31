package config

const (
	S3Bucket = "S3_BUCKET"
	Region   = "REGION"

	_defaultS3Bucket = "bucket"
	_defaultRegion   = "region"
)

type S3Config struct {
	S3Bucket string
	Region   string
}

func NewS3Config() S3Config {
	return S3Config{
		S3Bucket: GetString(S3Bucket, _defaultS3Bucket),
		Region:   GetString(Region, _defaultRegion),
	}
}
