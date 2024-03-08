package chat

import "time"

type Message struct {
	Type    int
	Message string    `json:"Message"`
	Time    time.Time `json:"Time"`
	Name    string    `json:"Name"`
}
