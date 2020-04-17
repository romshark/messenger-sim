package event

import "github.com/romshark/messenger-sim/messenger/eventlog"

type UserJoinedConversation struct {
	User         UserID
	Conversation ConversationID
}

// Copy creates a deep copy
func (e *UserJoinedConversation) Copy() eventlog.Payload {
	cp := *e
	return &cp
}
