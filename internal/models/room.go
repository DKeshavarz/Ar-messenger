package models

import "sync"

type Room struct {
	RoomName  string
	Clients   map[*Client]bool
	Broadcast chan Message
	Join      chan *Client
	Leave     chan *Client
	mu        sync.RWMutex 
}

func (r *Room) Run() {
    for {
        select {
        case client := <-r.Join:
            r.mu.Lock()
            r.Clients[client] = true
            r.mu.Unlock()

        case client := <-r.Leave:
            r.mu.Lock()
            delete(r.Clients, client)
            r.mu.Unlock()
            client.Conn.Close()
			
        case msg := <-r.Broadcast:
            r.mu.RLock()
            for client := range r.Clients {
                if client.RoomName == msg.RoomName {
                    if err := client.Conn.WriteJSON(msg); err != nil {
                        r.mu.RUnlock()
                        r.Leave <- client
                        r.mu.RLock()
                    }
                }
            }
            r.mu.RUnlock()
        }
    }
}