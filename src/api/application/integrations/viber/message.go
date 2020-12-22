package viber

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/application/integrations"
)

const (
	messageType = "text"
)

// MessageRequest is sent to Vibers Message API
type MessageRequest struct {
	Receiver     string   `json:"receiver"`
	Type         string   `json:"type"`
	Text         string   `json:"text,omitempty"`
	TrackingData string   `json:"tracking_data,omitempty"`
	Contact      *Contact `json:"contact,omitempty"`
}

// MessageResponse is returnned from Vibers Message API
type MessageResponse struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
	ChatHostname  string `json:"chat_hostname"`
}

// SendMessage sends message to Vibers API
func SendMessage(request *MessageRequest, viberBaseURL string, apiKey string) error {
	apiURL := fmt.Sprintf("%s/pa/send_message", viberBaseURL)
	authHeader := NewAuthHeader(apiKey)
	var response MessageResponse
	if err := integrations.SendPostRequest(apiURL, &request, &response, authHeader); err != nil {
		return err
	}
	if response.Status != SuccessStatus {
		return errors.Errorf("Request to %s failed with Status=%v and StatusMessage='%s'", apiURL, response.Status, response.StatusMessage)
	}
	return nil
}
