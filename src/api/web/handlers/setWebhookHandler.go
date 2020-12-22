package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/triviy/parklakes-viberbot/application/commands"
)

// SetWebhookHandler handles set webhook request
type SetWebhookHandler struct {
	cmd *commands.SetWebhookCmd
}

// NewSetWebhookHandler creates new handler instance
func NewSetWebhookHandler(cmd *commands.SetWebhookCmd) *SetWebhookHandler {
	return &SetWebhookHandler{cmd}
}

// Handle sets webhook url for Viber API callbacks
func (h SetWebhookHandler) Handle(c echo.Context) error {
	err := h.cmd.Execute()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, createOkResponse())
}
