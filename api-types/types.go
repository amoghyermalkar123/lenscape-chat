package apitypes

import (
	"time"
)

type ChannelCreation struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

type ChatMessage struct {
	Sender    string `json:"sender_email"`
	Message   string `json:"message"`
	CreatedAt time.Time
}
