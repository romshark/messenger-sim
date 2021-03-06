package event

import "github.com/romshark/messenger-sim/messenger/eventlog"

type MessageSent struct {
	ID           MessageID
	Body         string
	Sender       UserID
	Conversation ConversationID
}

// Copy creates a deep copy
func (e *MessageSent) Copy() eventlog.Payload {
	cp := *e
	return &cp
}
