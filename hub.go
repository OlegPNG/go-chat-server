package main

type Hub struct {
    clients map[*Client]bool
    broadcast chan []byte
    register chan *Client
    unregister chan *Client
}

func newHub() *Hub {
    return &Hub{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
    }
}

func (h *Hub) checkUserExists(username string) bool {
    for client := range h.clients {
	if client.username == username {
	    return true 
	}
    }

    return false
}

func (h *Hub) run() {
    for {
	select {
	case client := <- h.register:
	    if !h.checkUserExists(client.username) {
		h.clients[client] = true
	    } else {
		close(client.send)
	    }
	case client := <- h.unregister:
	    if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	    }
	case message := <- h.broadcast:
	    for client := range h.clients {
		select {
		case client.send <- message:
		default:
		    close(client.send)
		    delete(h.clients, client)
		}
	    }
	}
    }
}
