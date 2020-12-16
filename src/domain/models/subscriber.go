package models

// Subscriber of Viber bot
type Subscriber struct {
	ID           string   `bson:"_id"`
	Name         string   `bson:"name"`
	Avatar       string   `bson:"avatar,omitempty"`
	Country      string   `bson:"country,omitempty"`
	PhoneNumbers []string `bson:"phoneNumbers"`
	Active       bool     `bson:"active"`
}
