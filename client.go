package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
    // Time allowed to write message
    writeWait = 10 * time.Second

    // Time to read next pong message
    pongWait = 60 * time.Second

    // Must be less than pongWait
    pingPeriod = (pongWait * 9) / 10

    maxMessageSize = 512
)

var (
    newline = []byte{'\n'}
    space = []byte{' '}
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

type Client struct {
    username string
    hub *Hub
    conn *websocket.Conn
    send chan Message 
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil
    })

    for {
        _, raw, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        // Idk what this is doing
        //content := bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
        msg := Message{}
        err = json.Unmarshal(raw, &msg)
        if err != nil {
            log.Printf("Error umarshalling message: %v", err)
            errMsg := Message{
                Username: "server",
                Content: []byte("Error reading message"),
                Timestamp: time.Now(),
            }
            c.send <- errMsg
        } else {
            c.hub.broadcast <- msg
        }
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    for {
        select {
        case message, ok := <- c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // The hub closed the channel
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }

            raw, err := json.Marshal(message)
            if err != nil {
                log.Printf("Error marshalling message: %v", err)
                errMsg := ServerMessage([]byte("Could not send message"))
                w.Write(errMsg)
            } else {
                w.Write(raw)
            }

            n := len(c.send)
            for i := 0; i < n; i++ {
                raw, err := json.Marshal(<-c.send)
                if err != nil {
                    log.Printf("Error marshalling message: %v", err)
                    errMsg := ServerMessage([]byte("Could not send message"))
                    w.Write(errMsg)
                    break
                }
                w.Write(raw)
                //w.Write(newline)
                //w.Write(<-c.send)
            }

            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    username := r.Header.Get("username")
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    log.Printf("username: %v", username)
    client := &Client{username: username, hub: hub, conn: conn, send: make(chan Message, 256)}
    client.hub.register <- client


    go client.writePump()
    go client.readPump()
}

