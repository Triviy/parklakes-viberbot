package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/triviy/parklakes-viberbot/application/commands"
	"github.com/triviy/parklakes-viberbot/application/integrations/viber"
)

// CallbackHandler handles set webhook request
type CallbackHandler struct {
	getCarOwnerByTextCmd  *commands.GetCarOwnerByTextCmd
	getCarOwnerByImageCmd *commands.GetCarOwnerByImageCmd
	updateSubscriberCmd   *commands.UpdateSubscriberCmd
	unsubscribeCmd        *commands.UnsubscribeCmd
	welcomeCmd            *commands.WelcomeCmd
	inMemoryCache         *cache.Cache
}

// NewCallbackHandler creates new handler instance
func NewCallbackHandler(
	getCarOwnerByTextCmd *commands.GetCarOwnerByTextCmd,
	getCarOwnerByImageCmd *commands.GetCarOwnerByImageCmd,
	updateSubscriberCmd *commands.UpdateSubscriberCmd,
	unsubscribeCmd *commands.UnsubscribeCmd,
	welcomeCmd *commands.WelcomeCmd,
	inMemoryCache *cache.Cache,
) *CallbackHandler {
	return &CallbackHandler{
		getCarOwnerByTextCmd,
		getCarOwnerByImageCmd,
		updateSubscriberCmd,
		unsubscribeCmd,
		welcomeCmd,
		inMemoryCache,
	}
}

// Handle sets webhook url for Viber API callbacks
func (h CallbackHandler) Handle(c echo.Context) error {
	var r viber.Callback
	if err := c.Bind(&r); err != nil {
		return errors.Wrap(err, "binding of callback failed")
	}
	if _, ok := h.inMemoryCache.Get(string(r.MessageToken)); ok {
		return c.JSON(http.StatusOK, createOkResponse())
	}

	var res interface{} = createOkResponse()
	switch r.Event {
	case viber.SubscribedEvent:
		if err := h.updateSubscriberCmd.Execute(r.User, nil); err != nil {
			return err
		}
	case viber.UnsubscribedEvent:
		if err := h.unsubscribeCmd.Execute(r.UserID); err != nil {
			return err
		}
	case viber.ConversationStartedEvent:
		res = h.welcomeCmd.Execute()
	case viber.MessageEvent:
		var sendErr error
		if r.Message.Type == viber.PictureType {
			sendErr = h.getCarOwnerByImageCmd.Execute(r.Message, r.Sender.ID)
		} else {
			sendErr = h.getCarOwnerByTextCmd.Execute(r.Message, r.Sender.ID)
		}
		if updateErr := h.updateSubscriberCmd.Execute(r.Sender, r.Message.Contact); updateErr != nil {
			log.Error(updateErr)
		}
		if sendErr != nil {
			return sendErr
		}
	}

	h.inMemoryCache.SetDefault(string(r.MessageToken), true)
	return c.JSON(http.StatusOK, res)
}
