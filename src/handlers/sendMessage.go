package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/triviy/parklakes-viberbot/config"
)

const (
	messageType = "text"
)

// {
// 	"receiver":"01234567890A=",
// 	"min_api_version":1,
// 	"sender":{
// 	   "name":"John McClane",
// 	   "avatar":"http://avatar.example.com"
// 	},
// 	"tracking_data":"tracking data",
// 	"type":"text",
// 	"text":"Hello world!"
//  }
type sendMessageRequest struct {
	AuthToken string `json:"auth_token"`
	Receiver  string `json:"receiver"`
	Type      string `json:"type"`
	Text      string `json:"text"`
}

// {
// 	"status":0,
// 	"status_message":"ok",
// 	"chat_hostname":"data",
//  }
type sendMessageResponse struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
	ChatHostname  string `json:"chat_hostname"`
}

// SendMessage sends Viber user a message
func SendMessage(receiver string, message string) {
	request := sendMessageRequest{
		AuthToken: config.GetViberAPIKey(),
		Receiver:  receiver,
		Type:      messageType,
		Text:      message,
	}

	bytesRepresentation, err := json.Marshal(request)
	if err != nil {
		log.Fatalln(err)
	}

	apiURL := fmt.Sprintf("%s/pa/send_message", config.GetViberBaseURL())
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result sendMessageResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(result)
}
