package viber

const (
	// ViberSuccessStatus represents success status for any Viber API calls
	ViberSuccessStatus = 0
)

type Callback struct {
	Event        string          `json:"event"`
	Timestamp    int64           `json:"timestamp"`
	MessageToken int64           `json:"message_token"`
	UserID       string          `json:"user_id,omitempty"`
	Type         string          `json:"type,omitempty"`
	Context      string          `json:"context,omitempty"`
	Subscribed   bool            `json:"subscribed,omitempty"`
	User         User            `json:"user,omitempty"`
	Message      CallbackMessage `json:"message,omitempty"`
}

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar,omitempty"`
	Country    string `json:"country,omitempty"`
	Language   string `json:"language,omitempty"`
	APIVersion int    `json:"api_version,omitempty"`
}

type CallbackMessage struct {
	Type         string   `json:"type"`
	Text         string   `json:"text"`
	Media        string   `json:"media,omitempty"`
	TrackingData string   `json:"tracking_data,omitempty"`
	FileName     string   `json:"file_name,omitempty"`
	FileSize     int64    `json:"file_size,omitempty"`
	Duration     int64    `json:"duration,omitempty"`
	StickerID    string   `json:"sticker_id,omitempty"`
	Location     Location `json:"location,omitempty"`
	Contact      Contact  `json:"contact,omitempty"`
}

type Location struct {
	Lat float64 `json:"id"`
	Lon float64 `json:"lon"`
}

type Contact struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Avatar      string `json:"avatar,omitempty"`
}

type WelcomeResponse struct {
	TrackingData string `json:"tracking_data,omitempty"`
	Type         string `json:"type"`
	Text         string `json:"text"`
}

// Subscribed
// {
// 	"event":"subscribed",
// 	"timestamp":1457764197627,
// 	"user":{
// 	   "id":"01234567890A=",
// 	   "name":"John McClane",
// 	   "avatar":"http://avatar.example.com",
// 	   "country":"UK",
// 	   "language":"en",
// 	   "api_version":1
// 	},
// 	"message_token":4912661846655238145
//  }

// Welcome response:
// {
// 	"sender":{
// 	   "name":"John McClane",
// 	   "avatar":"http://avatar.example.com"
// 	},
// 	"tracking_data":"tracking data",
// 	"type":"picture",
// 	"text":"Welcome to our bot!",
// 	"media":"http://www.images.com/img.jpg",
// 	"thumbnail":"http://www.images.com/thumb.jpg"
//  }

// Unsubscribed
// {
// 	"event":"unsubscribed",
// 	"timestamp":1457764197627,
// 	-->"user_id":"01234567890A=",
// 	"message_token":4912661846655238145
//  }

// Conversation started
// {
// 	"event":"conversation_started",
// 	"timestamp":1457764197627,
// 	"message_token":4912661846655238145,
// 	"type":"open",
// 	"context":"context information",
// 	"user":{
// 	   "id":"01234567890A=",
// 	   "name":"John McClane",
// 	   "avatar":"http://avatar.example.com",
// 	   "country":"UK",
// 	   "language":"en",
// 	   "api_version":1
// 	},
// 	"subscribed":false
//  }

// Message
// {
// 	"event":"message",
// 	"timestamp":1457764197627,
// 	"message_token":4912661846655238145,
// 	"sender":{
// 	   "id":"01234567890A=",
// 	   "name":"John McClane",
// 	   "avatar":"http://avatar.example.com",
// 	   "country":"UK",
// 	   "language":"en",
// 	   "api_version":1
// 	},
// 	"message":{
// 	   "type":"text",
// 	   "text":"a message to the service",
// 	   "media":"http://example.com",
// 	   "location":{
// 		  "lat":50.76891,
// 		  "lon":6.11499
// 	   },
// 	   "tracking_data":"tracking data"
//		"contact":"" {
//			"name":""
//			"phone_number":""
//			"avatar":""
//		}
//		"tracking_data":""
//		"file_name":""
//		"file_size":""
//		"duration":""
//		"sticker_id":"URL"
// 	}
//  }
