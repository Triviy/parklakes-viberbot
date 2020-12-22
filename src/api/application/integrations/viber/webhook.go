package viber

// SetWebhookRequest request for Vibers webhook API
type SetWebhookRequest struct {
	URL        string   `json:"url"`
	EventTypes []string `json:"event_types"`
	SendName   bool     `json:"send_name"`
	SendPhoto  bool     `json:"send_photo"`
}

// SetWebhookResponse response from Viber webhook API
type SetWebhookResponse struct {
	Status        int      `json:"status"`
	StatusMessage string   `json:"status_message"`
	EventTypes    []string `json:"event_types"`
}
