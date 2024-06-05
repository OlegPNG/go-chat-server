package main

type Hub struct {
    clients map[*Client]bool
    broadcast chan Message 
    register chan *Client
    unregister chan *Client
    history []Message
}

func newHub() *Hub {
    return  &Hub{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan Message),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	//history: make([]Message, 0),
	history: testHistory(),
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
	    h.history = append(h.history, message)
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
