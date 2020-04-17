package model

import (
	"time"

	"github.com/romshark/messenger-sim/messenger/event"
)

type Conversation struct {
	ConversationID event.ConversationID

	ID           string    `json:"id"`
	Title        string    `json:"title"`
	AvatarURL    *string   `json:"avatarURL"`
	CreationTime time.Time `json:"creationTime"`
}
