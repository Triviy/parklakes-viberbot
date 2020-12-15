package viber

// Callback is Viber full callback contract
type Callback struct {
	Event        string          `json:"event"`
	Timestamp    int64           `json:"timestamp"`
	MessageToken int64           `json:"message_token"`
	UserID       string          `json:"user_id,omitempty"`
	Type         string          `json:"type,omitempty"`
	Context      string          `json:"context,omitempty"`
	Subscribed   bool            `json:"subscribed,omitempty"`
	Sender       User            `json:"sender,omitempty"`
	User         User            `json:"user,omitempty"`
	Message      CallbackMessage `json:"message,omitempty"`
}

// User is Viber full users contract
type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar,omitempty"`
	Country    string `json:"country,omitempty"`
	Language   string `json:"language,omitempty"`
	APIVersion int    `json:"api_version,omitempty"`
}

// CallbackMessage is Viber message details
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

// Location is Viber users location
type Location struct {
	Lat float64 `json:"id"`
	Lon float64 `json:"lon"`
}

// Contact is Viber users contacts
type Contact struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Avatar      string `json:"avatar,omitempty"`
}
