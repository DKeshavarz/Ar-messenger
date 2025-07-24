package models

type Room struct {
	RoomName  string
	Clients   map[*Client]bool
	Broadcast chan Message
	Join      chan *Client
	Leave     chan *Client
	//mu        sync.RWMutex -> TODO: implement mutex for thread safety
}
