package event

import "github.com/romshark/messenger-sim/messenger/eventlog"

type UserLeftConversation struct {
	User         UserID
	Conversation ConversationID
}

// Copy creates a deep copy
func (e *UserLeftConversation) Copy() eventlog.Payload {
	cp := *e
	return &cp
}
