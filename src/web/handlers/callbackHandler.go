package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/triviy/parklakes-viberbot/application/commands"
	"github.com/triviy/parklakes-viberbot/application/integrations/viber"
)

// CallbackHandler handles set webhook request
type CallbackHandler struct {
	cmd *commands.SetWebhookCmd
}

// NewCallbackHandler creates new handler instance
func NewCallbackHandler(cmd *commands.SetWebhookCmd) *CallbackHandler {
	return &CallbackHandler{cmd}
}

// Handle sets webhook url for Viber API callbacks
func (h CallbackHandler) Handle(c echo.Context) error {
	var r viber.Callback
	if err := c.Bind(&r); err != nil {
		return errors.Wrap(err, "binding of callback failed")
	}
	switch r.Event {
	case viber.SubscribedEvent:
		logrus.Info("Save data about user")
	case viber.UnsubscribedEvent:
		logrus.Info("Remove from database")
	case viber.ConversationStartedEvent:
		logrus.Info("Response with Welcome request")
	case viber.MessageEvent:
		logrus.Info("Send message to user")
		logrus.Info("Save data about user")
		logrus.Info("Save contacts")
	default:
		logrus.Info("Log")
	}

	err := h.cmd.Execute()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, createOkResponse())
}
