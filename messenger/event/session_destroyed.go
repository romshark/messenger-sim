package event

import (
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/sessid"
)

type SessionDestroyed struct {
	Session sessid.SessionID
}

// Copy creates a deep copy
func (e *SessionDestroyed) Copy() eventlog.Payload {
	cp := *e
	return &cp
}
