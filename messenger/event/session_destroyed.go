package event

import (
	"simulator/messenger/eventlog"
	"simulator/messenger/sessid"
)

type SessionDestroyed struct {
	Session sessid.SessionID
}

// Copy creates a deep copy
func (e *SessionDestroyed) Copy() eventlog.Payload {
	cp := *e
	return &cp
}
