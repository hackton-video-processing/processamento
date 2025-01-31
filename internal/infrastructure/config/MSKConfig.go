package config

const (
	TopicName = "TOPIC_NAME"

	_defaultTopicName = "test"
)

type MSKConfig struct {
	TopicName string
}

func NewMSKConfig() MSKConfig {
	return MSKConfig{
		TopicName: GetString(TopicName, _defaultTopicName),
	}
}
