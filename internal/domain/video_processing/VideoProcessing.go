package video_processing

type VideoProcessing struct {
	ExecutionID string `json:"execution_id"`
	VideoName   string `json:"video_name"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}
