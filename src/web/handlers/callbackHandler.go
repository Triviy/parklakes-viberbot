package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/triviy/parklakes-viberbot/application/commands"
	"github.com/triviy/parklakes-viberbot/application/integrations/viber"
)

// CallbackHandler handles set webhook request
type CallbackHandler struct {
	getCarOwnerByTextCmd *commands.GetCarOwnerByTextCmd
	updateSubscriberCmd  *commands.UpdateSubscriberCmd
	unsubscribeCmd       *commands.UnsubscribeCmd
	welcomeCmd           *commands.WelcomeCmd
}

// NewCallbackHandler creates new handler instance
func NewCallbackHandler(
	getCarOwnerByTextCmd *commands.GetCarOwnerByTextCmd,
	updateSubscriberCmd *commands.UpdateSubscriberCmd,
	unsubscribeCmd *commands.UnsubscribeCmd,
	welcomeCmd *commands.WelcomeCmd,
) *CallbackHandler {
	return &CallbackHandler{
		getCarOwnerByTextCmd,
		updateSubscriberCmd,
		unsubscribeCmd,
		welcomeCmd,
	}
}

// Handle sets webhook url for Viber API callbacks
func (h CallbackHandler) Handle(c echo.Context) error {
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return errors.Unwrap(err)
	}

	var r viber.Callback
	if err := json.Unmarshal(b, &r); err != nil {
		return errors.Wrap(err, "binding of callback failed")
	}

	switch r.Event {
	case viber.SubscribedEvent:
		if err := h.updateSubscriberCmd.Execute(&r.User, nil); err != nil {
			return err
		}
	case viber.UnsubscribedEvent:
		if err := h.unsubscribeCmd.Execute(r.UserID); err != nil {
			return err
		}
	case viber.ConversationStartedEvent:
		r := h.welcomeCmd.Execute()
		return c.JSON(http.StatusOK, r)
	case viber.MessageEvent:
		sendErr := h.getCarOwnerByTextCmd.Execute(r.Message.Text, r.Sender.ID, r.Message.TrackingData)
		if updateErr := h.updateSubscriberCmd.Execute(&r.User, &r.Message.Contact); updateErr != nil {
			logrus.Error(updateErr)
		}
		if sendErr != nil {
			return sendErr
		}
	}

	return c.JSON(http.StatusOK, createOkResponse())
}
