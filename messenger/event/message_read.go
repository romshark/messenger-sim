package event

import "github.com/romshark/messenger-sim/messenger/eventlog"

type MessageRead struct {
	Message MessageID
	User    UserID
}

// Copy creates a deep copy
func (e *MessageRead) Copy() eventlog.Payload {
	cp := *e
	return &cp
}
