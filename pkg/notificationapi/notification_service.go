package notificationapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type NotificationService struct {
	Endpoint string
}

func NewNotificationService(baseURL, endpoint string) (*NotificationService, error) {
	return &NotificationService{Endpoint: fmt.Sprintf("%s/%s", baseURL, endpoint)}, nil
}

type NotificationRequest struct {
	Email   string `json:"userEmail"`
	Message string `json:"message"`
}

func (c *NotificationService) SendNotification(email, message string) error {
	requestBody := NotificationRequest{
		Email:   email,
		Message: message,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := http.Post(c.Endpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
