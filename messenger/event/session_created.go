package event

import (
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/sessid"
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
