package config

const (
	BaseURL        = "NOTIFICATION_API_BASE_URL"
	NotifyEndpoint = "NOTIFICATION_API_ENDPOINT"

	_defaultBaseURL        = "http://18.234.220.7:80"
	_defaultNotifyEndpoint = "v1/notification"
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
