package models

import "github.com/gorilla/websocket"

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	RoomName string `json:"room_name"`
	Username string `json:"username"`
}


