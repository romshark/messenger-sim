package model

import (
	"simulator/messenger/event"
	"time"
)

type Conversation struct {
	ConversationID event.ConversationID

	ID           string    `json:"id"`
	Title        string    `json:"title"`
	AvatarURL    *string   `json:"avatarURL"`
	CreationTime time.Time `json:"creationTime"`
}
