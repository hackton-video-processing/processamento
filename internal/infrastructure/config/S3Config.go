package config

const (
	S3Bucket        = "S3_BUCKET"
	AccessKey       = "AWS_ACCESS_KEY_ID"
	SecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	Region          = "AWS_REGION"
	DownloadPath    = "DOWNLOAD_PATH"
	UploadPath      = "UPLOAD_PATH"

	_defaultS3Bucket        = "video-processing-api-bucket"
	_defaultRegion          = "us-east-1"
	_defaultAccessKey       = "access_key"
	_defaultSecretAccessKey = "secret_accessKey"
	_defaultDownloadPath    = "upload/"
	_defaultUploadPath      = "download/"
)

type S3Config struct {
	S3Bucket        string
	Region          string
	AccessKey       string
	SecretAccessKey string
	DownloadPath    string
	UploadPath      string
}

func NewS3Config() S3Config {
	return S3Config{
		S3Bucket:        GetString(S3Bucket, _defaultS3Bucket),
		Region:          GetString(Region, _defaultRegion),
		AccessKey:       GetString(AccessKey, _defaultAccessKey),
		SecretAccessKey: GetString(SecretAccessKey, _defaultSecretAccessKey),
		DownloadPath:    GetString(DownloadPath, _defaultDownloadPath),
		UploadPath:      GetString(UploadPath, _defaultUploadPath),
	}
}
