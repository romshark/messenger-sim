package event

import "github.com/romshark/messenger-sim/messenger/eventlog"

type UserRemovedFromConversation struct {
	Conversation ConversationID
	Remover      UserID
	Removed      UserID
	Reason       *string
}

// Copy creates a deep copy
func (e *UserRemovedFromConversation) Copy() eventlog.Payload {
	cp := *e

	if e.Reason != nil {
		v := *e.Reason
		cp.Reason = &v
	}

	return &cp
}
