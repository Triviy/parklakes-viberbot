package handlers

import (
	"fmt"
	"net/http"
	"time"

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

	messageID := fmt.Sprint(r.MessageToken)
	if processed := h.lockCallback(messageID); processed {
		log.Infof("Message with token %s was already processed", messageID)
		return c.JSON(http.StatusOK, createOkResponse())
	}

	for i := 0; i < 5; i++ {
		if v, ok := h.inMemoryCache.Get(messageID); ok {
			if v.(bool) {
				log.Infof("Message with token %s was already processed", messageID)
				return c.JSON(http.StatusOK, createOkResponse())
			}
			time.Sleep(time.Second)
		} else {
			h.inMemoryCache.SetDefault(messageID, false)
			break
		}
	}

	res, err := h.handleCallback(r)
	if err != nil {
		return err
	}

	h.unlockCallback(messageID)
	return c.JSON(http.StatusOK, res)
}

func (h CallbackHandler) lockCallback(messageID string) (processed bool) {
	log.Info(")))) start")
	for i := 0; i < 5; i++ {
		if v, ok := h.inMemoryCache.Get(messageID); ok {
			log.Info(")))) ok is true")
			if v.(bool) {
				log.Info(")))) v is true")
				return true
			}
			log.Info(")))) v is false, sleeping")
			time.Sleep(time.Second)
		} else {
			log.Info(")))) ok is false")
			h.inMemoryCache.SetDefault(messageID, false)
			return false
		}
	}
	return false
}

func (h CallbackHandler) unlockCallback(messageID string) {
	h.inMemoryCache.SetDefault(messageID, true)
}

func (h CallbackHandler) handleCallback(r viber.Callback) (res interface{}, err error) {
	res = createOkResponse()
	switch r.Event {
	case viber.SubscribedEvent:
		if err := h.updateSubscriberCmd.Execute(r.User, nil); err != nil {
			return res, err
		}
	case viber.UnsubscribedEvent:
		if err := h.unsubscribeCmd.Execute(r.UserID); err != nil {
			return res, err
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
			return res, sendErr
		}
	}
	return
}
