package event

import "simulator/messenger/eventlog"

type UserLeftConversation struct {
	User         UserID
	Conversation ConversationID
}

// Copy creates a deep copy
func (e *UserLeftConversation) Copy() eventlog.Payload {
	cp := *e
	return &cp
}
