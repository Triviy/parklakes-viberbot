package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/triviy/parklakes-viberbot/config"
)

// {
// 	"url":"https://my.host.com",
// 	"event_types":[
// 	   "delivered",
// 	   "seen",
// 	   "failed",
// 	   "subscribed",
// 	   "unsubscribed",
// 	   "conversation_started"
// 	],
// 	"send_name": true,
// 	"send_photo": true
//  }
type setWebhookRequest struct {
	AuthToken  string   `json:"auth_token"`
	URL        string   `json:"url"`
	EventTypes []string `json:"event_types"`
	SendName   bool     `json:"send_name"`
	SendPhoto  bool     `json:"send_photo"`
}

// {
// 	"status":0,
// 	"status_message":"ok",
// 	"event_types":[
// 	   "delivered",
// 	   "seen",
// 	   "failed",
// 	   "subscribed",
// 	   "unsubscribed",
// 	   "conversation_started"
// 	]
//  }
type setWebhookResponse struct {
	Status        int      `json:"status"`
	StatusMessage string   `json:"status_message"`
	EventTypes    []string `json:"event_types"`
}

// SetWebhook sends a Webhook url for Viber API
func SetWebhook(c echo.Context) error {
	request := setWebhookRequest{
		AuthToken: config.GetViberAPIKey(),
		URL:       config.GetViberWebhookURL(),
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
		log.Fatalln(err)
	}
	apiURL := fmt.Sprintf("%s/pa/set_webhook", config.GetViberBaseURL())
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result setWebhookResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(result)
}
