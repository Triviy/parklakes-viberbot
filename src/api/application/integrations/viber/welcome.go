package viber

// WelcomeResponse is a response to conversation_started callback event
type WelcomeResponse struct {
	TrackingData string `json:"tracking_data,omitempty"`
	Type         string `json:"type"`
	Text         string `json:"text"`
}
