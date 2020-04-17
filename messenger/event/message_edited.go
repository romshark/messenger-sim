package event

import "simulator/messenger/eventlog"

type MessageEdited struct {
	Message MessageID
	Editor  UserID
	Body    string
}

// Copy creates a deep copy
func (e *MessageEdited) Copy() eventlog.Payload {
	cp := *e
	return &cp
}
