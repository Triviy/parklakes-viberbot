package viber

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/application/integrations"
)

const (
	messageType = "text"
)

type MessageRequest struct {
	//TODO: Move to headers
	AuthToken    string  `json:"auth_token"`
	Receiver     string  `json:"receiver"`
	Type         string  `json:"type"`
	Text         string  `json:"text,omitempty"`
	TrackingData string  `json:"tracking_data,omitempty"`
	Contact      Contact `json:"contact,omitempty"`
}

type MessageResponse struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
	ChatHostname  string `json:"chat_hostname"`
}

func SendMessage(request *MessageRequest, viberBaseURL string) error {
	apiURL := fmt.Sprintf("%s/pa/send_message", viberBaseURL)

	var response MessageResponse
	if err := integrations.SendPostRequest(apiURL, &request, &response); err != nil {
		return err
	}
	if response.Status != SuccessStatus {
		return errors.Errorf("Request to %s failed with Status=%v and StatusMessage='%s'", apiURL, response.Status, response.StatusMessage)
	}
	return nil
}
