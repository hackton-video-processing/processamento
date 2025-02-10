package config

const (
	BaseURL        = "NOTIFICATION_API_BASE_URL"
	NotifyEndpoint = "NOTIFICATION_API_ENDPOINT"

	_defaultBaseURL        = "url"
	_defaultNotifyEndpoint = "endpoint"
)

type NotificationAPIConfig struct {
	BaseURL  string
	Endpoint string
}

func NewNotificationAPIConfig() NotificationAPIConfig {
	return NotificationAPIConfig{
		BaseURL:  GetString(BaseURL, _defaultBaseURL),
		Endpoint: GetString(NotifyEndpoint, _defaultNotifyEndpoint),
	}
}
