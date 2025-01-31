package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	Env                   Environment
	Port                  string
	KafkaConfig           MSKConfig
	S3Config              S3Config
	VideoProcessingConfig VideoProcessingConfig
}

func LoadConfiguration() (AppConfig, error) {
	return AppConfig{
		Env:                   GetEnvironment(),
		Port:                  GetPort(),
		KafkaConfig:           NewMSKConfig(),
		S3Config:              NewS3Config(),
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

func GetBool(env string, defaultValue bool) bool {
	value := os.Getenv(env)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

func GetEnvironment() Environment {
	value := os.Getenv("ENVIRONMENT")
	if value == "" {
		fmt.Println(fmt.Sprintf("%s: %s", "ENVIRONMENT", "local"))
		return LOCAL
	}

	if value == string(PROD) {
		fmt.Println(fmt.Sprintf("%s: %s", "ENVIRONMENT", "production"))
		return PROD
	}
	fmt.Println(fmt.Sprintf("%s: %s", "ENVIRONMENT", "local"))
	return LOCAL
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
