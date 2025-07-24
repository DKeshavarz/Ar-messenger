package models

type Message struct {
	Content  string `json:"content"`
	RoomName string `json:"room_name"`
	Username string `json:"username"`
}