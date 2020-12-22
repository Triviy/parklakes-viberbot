package viber

import "github.com/triviy/parklakes-viberbot/application/integrations"

// Viber constants
const (
	SuccessStatus = 0

	APIKeyHeader = "X-Viber-Auth-Token"

	DeliveredEvent           = "delivered"
	SeenEvent                = "seen"
	FailedEvent              = "failed"
	SubscribedEvent          = "subscribed"
	UnsubscribedEvent        = "unsubscribed"
	ConversationStartedEvent = "conversation_started"
	MessageEvent             = "message"

	TextType    = "text"
	PictureType = "picture"
)

// NewAuthHeader creates Viber auth header
func NewAuthHeader(apiKey string) integrations.Header {
	return integrations.Header{Name: APIKeyHeader, Value: apiKey}
}
