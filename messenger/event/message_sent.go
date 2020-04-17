package event

import "simulator/messenger/eventlog"

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
