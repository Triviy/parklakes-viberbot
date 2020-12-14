package viber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const (
	messageType = "text"
)

type MessageRequest struct {
	AuthToken    string  `json:"auth_token"`
	Receiver     string  `json:"receiver"`
	Type         string  `json:"type"`
	Text         string  `json:"text",omitempty`
	TrackingData string  `json:"tracking_data,omitempty"`
	Contact      Contact `json:"contact,omitempty"`
}

type messageResponse struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
	ChatHostname  string `json:"chat_hostname"`
}

func sendViberMessage(request *MessageRequest, viberBaseURL string) error {
	bytesRepresentation, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "serialization of message request failed")
	}

	apiURL := fmt.Sprintf("%s/pa/send_message", viberBaseURL)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return errors.Wrapf(err, "sending request to %s failed", apiURL)
	}

	var result messageResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return errors.Wrap(err, "deserialization of response failed")
	}
	if result.Status != ViberSuccessStatus {
		return errors.Errorf("Request to %s failed with Status=%v and StatusMessage='%s'", apiURL, result.Status, result.StatusMessage)
	}
	return nil
}
