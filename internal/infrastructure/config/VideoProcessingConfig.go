package config

const (
	MAX_VIDEO_PROCESSING = "MAX_VIDEO_PROCESSING"

	_defaultMaxVideoProcessing = 2
)

type VideoProcessingConfig struct {
	MaxVideoProcessing int
}

func NewVideoProcessingConfig() VideoProcessingConfig {
	return VideoProcessingConfig{
		MaxVideoProcessing: GetInt(MAX_VIDEO_PROCESSING, _defaultMaxVideoProcessing),
	}
}
