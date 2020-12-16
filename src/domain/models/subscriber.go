package models

// Subscriber of Viber bot
type Subscriber struct {
	ID           string   `json:"_id"`
	Name         string   `json:"name"`
	Avatar       string   `json:"avatar,omitempty"`
	Country      string   `json:"country,omitempty"`
	PhoneNumbers []string `json:"phoneNumbers"`
	Active       bool     `json:"active"`
}
