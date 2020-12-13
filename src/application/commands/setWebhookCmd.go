package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/application/integrations"
	"github.com/triviy/parklakes-viberbot/web/config"
)

// SetWebhookCmd runs migration
type SetWebhookCmd struct {
	config *config.APIConfig
}

// SetWebhookRequest request format to /pa/set_webhook
type SetWebhookRequest struct {
	AuthToken  string   `json:"auth_token"`
	URL        string   `json:"url"`
	EventTypes []string `json:"event_types"`
	SendName   bool     `json:"send_name"`
	SendPhoto  bool     `json:"send_photo"`
}

// SetWebhookResponse response format to /pa/set_webhook
type SetWebhookResponse struct {
	Status        int      `json:"status"`
	StatusMessage string   `json:"status_message"`
	EventTypes    []string `json:"event_types"`
}

// NewSetWebhookCmd creates new instance of SetWebhookCmd
func NewSetWebhookCmd(config *config.APIConfig) *SetWebhookCmd {
	return &SetWebhookCmd{config}
}

// Execute calls setting Viber callback URLs
func (cmd SetWebhookCmd) Execute() error {
	request := SetWebhookRequest{
		AuthToken: cmd.config.GetViberAPIKey(),
		URL:       cmd.config.GetViberWebhookURL(),
		EventTypes: []string{
			"delivered",
			"seen",
			"failed",
			"subscribed",
			"unsubscribed",
			"conversation_started",
		},
		SendName:  true,
		SendPhoto: true,
	}

	bytesRepresentation, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "request serialization to JSON failed")
	}
	apiURL := fmt.Sprintf("%s/pa/set_webhook", cmd.config.GetViberBaseURL())
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return errors.Wrapf(err, "calling %s failed", apiURL)
	}

	var result SetWebhookResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return errors.Wrapf(err, "response deserialization to JSON failed", apiURL)
	}
	log.Println(result)
	if result.Status != integrations.ViberSuccessStatus {
		return errors.Errorf("Request to %s failed with Status=%v and StatusMessage='%s'", apiURL, result.Status, result.StatusMessage)
	}
	return nil
}
