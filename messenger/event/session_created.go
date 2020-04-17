package event

import (
	"simulator/messenger/eventlog"
	"simulator/messenger/sessid"
)

type SessionCreated struct {
	ID        sessid.SessionID
	User      UserID
	IP        string
	UserAgent string
}

// Copy creates a deep copy
func (e *SessionCreated) Copy() eventlog.Payload {
	cp := *e
	return &cp
}
