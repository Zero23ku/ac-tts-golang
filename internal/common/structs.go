package common

type Condition struct {
	BroadcasterUserId string `json:"broadcaster_user_id"`
	UserId            string `json:"user_id"`
}

type Transport struct {
	Method    string `json:"method"`
	Callback  string `json:"callback"`
	Secret    string `json:"secret"`
	SessionId string `json:"session_id"`
}

type ChatSubscriptionEvent struct {
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Condition Condition `json:"condition"`
	Transport Transport `json:"transport"`
}

type Broadcaster struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadCasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	CreatedAt       string `json:"created_at"`
}

type Response struct {
	Data []Broadcaster `json:"data"`
}

type Fragments struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Message struct {
	Text      string      `json:"text"`
	Fragments []Fragments `json:"fragments"`
}

type Event struct {
	BroadcasterUserId    string  `json:"broadcaster_user_id"`
	BroadcasterUserLogin string  `json:"broadcaster_user_login"`
	BroadcasterUserName  string  `json:"broadcaster_user_name"`
	MessageId            string  `json:"message_id"`
	Message              Message `json:"message"`
}
