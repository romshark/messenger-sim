package model

import (
	"time"

	"github.com/romshark/messenger-sim/messenger/event"
)

type Message struct {
	MessageID      event.MessageID
	SenderID       event.UserID
	ConversationID event.ConversationID

	ID          string    `json:"id"`
	Body        string    `json:"body"`
	SendingTime time.Time `json:"sendingTime"`
}
