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

func testMessage() Message {
    return Message{
        Username: "server",
        Content: "This message is a test",
    }
}

func testHistory() []Message {
    hist := make([]Message, 0)
    hist = append(hist, testMessage())
    hist = append(hist, testMessage())
    hist = append(hist, testMessage())

    return hist
}
