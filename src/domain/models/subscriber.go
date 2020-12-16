package models

// Subscriber of Viber bot
type Subscriber struct {
	ID           string   `bson:"_id,omitempty"`
	Name         string   `bson:"name"`
	Avatar       string   `bson:"avatar,omitempty"`
	Country      string   `bson:"country,omitempty"`
	PhoneNumbers []string `bson:"phoneNumbers,omitempty"`
	Active       bool     `bson:"active"`
}
