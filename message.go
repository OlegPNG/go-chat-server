package main

import (
	"encoding/json"
	"time"
	//"encoding/json"
)

type Message struct {
    Username    string      `json:"username"`
    Content     []byte      `json:"content"`
    Timestamp   time.Time   `json:"timestamp"`
}

func ServerMessage(content []byte) []byte {
    msg, _ := json.Marshal(Message{
        Username: "server",
        Content: content,
        Timestamp: time.Now(),
    })

    return msg
}
