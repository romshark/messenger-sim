package event

import "simulator/messenger/eventlog"

type UserJoinedConversation struct {
	User         UserID
	Conversation ConversationID
}

// Copy creates a deep copy
func (e *UserJoinedConversation) Copy() eventlog.Payload {
	cp := *e
	return &cp
}
