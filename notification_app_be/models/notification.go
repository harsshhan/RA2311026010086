package models

type Notification struct {
	ID        string `json:"ID"`
	Type      string `json:"Type"`
	Message   string `json:"Message"`
	Timestamp string `json:"Timestamp"`
}

type NotificationsResponse struct {
	Notifications []Notification `json:"notifications"`
}
