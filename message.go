package main

import (
	"encoding/json"
)

type Message struct {
    Username    string      `json:"username"`
    Content     string      `json:"content"`
}

func ServerMessage(content string) []byte {
    msg, _ := json.Marshal(Message{
        Username: "server",
        Content: content,
    })

    return msg
}
