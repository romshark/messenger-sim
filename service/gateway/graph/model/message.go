package model

import (
	"simulator/messenger/event"
	"time"
)

type Message struct {
	MessageID      event.MessageID
	SenderID       event.UserID
	ConversationID event.ConversationID

	ID          string    `json:"id"`
	Body        string    `json:"body"`
	SendingTime time.Time `json:"sendingTime"`
}
