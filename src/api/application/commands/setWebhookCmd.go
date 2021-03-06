package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/application/integrations"
	"github.com/triviy/parklakes-viberbot/application/integrations/viber"
	"github.com/triviy/parklakes-viberbot/web/config"
)

// SetWebhookCmd instance of viber webhook cmd
type SetWebhookCmd struct {
	config *config.APIConfig
}

// NewSetWebhookCmd creates new instance of SetWebhookCmd
func NewSetWebhookCmd(config *config.APIConfig) *SetWebhookCmd {
	return &SetWebhookCmd{config}
}

// Execute calls setting Viber callback URLs
func (cmd SetWebhookCmd) Execute() error {
	webhookURL := fmt.Sprintf("https://%s/api/v1/viber/callback", cmd.config.GetAppBaseURL())
	request := viber.SetWebhookRequest{
		URL: webhookURL,
		EventTypes: []string{
			viber.DeliveredEvent,
			viber.SeenEvent,
			viber.FailedEvent,
			viber.SubscribedEvent,
			viber.UnsubscribedEvent,
			viber.ConversationStartedEvent,
		},
		SendName:  true,
		SendPhoto: true,
	}

	apiURL := fmt.Sprintf("%s/pa/set_webhook", cmd.config.GetViberBaseURL())
	authHeader := viber.NewAuthHeader(cmd.config.GetViberAPIKey())

	var response viber.SetWebhookResponse
	if err := integrations.SendPostRequest(apiURL, &request, &response, authHeader); err != nil {
		return err
	}
	if response.Status != viber.SuccessStatus {
		return errors.Errorf("Request to %s failed with Status=%v and StatusMessage='%s'", apiURL, response.Status, response.StatusMessage)
	}
	return nil
}
