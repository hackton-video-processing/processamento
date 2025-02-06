package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	Port                  string
	KafkaConfig           MSKConfig
	S3Config              S3Config
	MySQL                 MySQLConfig
	VideoProcessingConfig VideoProcessingConfig
}

func LoadConfiguration() (AppConfig, error) {
	return AppConfig{
		Port:                  GetPort(),
		KafkaConfig:           NewMSKConfig(),
		S3Config:              NewS3Config(),
		MySQL:                 NewMySQLConfig(),
		VideoProcessingConfig: NewVideoProcessingConfig(),
	}, nil
}

func GetString(env string, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		fmt.Println(fmt.Sprintf("%s: %s", env, defaultValue))
		return defaultValue
	}
	fmt.Println(fmt.Sprintf("%s: %s", env, value))
	return value
}

func GetInt(env string, defaultValue int) int {
	value := os.Getenv(env)
	if value == "" {
		fmt.Println(fmt.Sprintf("%s: %d", env, defaultValue))
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s: %d", env, defaultValue))
		return defaultValue
	}
	fmt.Println(fmt.Sprintf("%s: %d", env, intValue))
	return intValue
}

func GetPort() string {
	value := os.Getenv("PORT")
	if value == "" {
		fmt.Println(fmt.Sprintf("%s: %s", "PORT", "8080"))
		return "8080"
	}
	fmt.Println(fmt.Sprintf("%s: %s", "PORT", value))
	return value
}
