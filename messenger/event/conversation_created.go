package event

import (
	"net/url"
	"simulator/messenger/eventlog"
)

type ConversationCreated struct {
	ID           ConversationID
	Creator      UserID
	Participants []UserID
	Title        string
	AvatarURL    *url.URL
}

// Copy creates a deep copy
func (e *ConversationCreated) Copy() eventlog.Payload {
	cp := *e

	if e.AvatarURL != nil {
		v := *e.AvatarURL
		cp.AvatarURL = &v
	}

	if e.Participants != nil {
		cp.Participants = make([]UserID, len(e.Participants))
		copy(cp.Participants, e.Participants)
	}

	return &cp
}
