package event

import "github.com/romshark/messenger-sim/messenger/eventlog"

type MessageDeleted struct {
	Message MessageID
	Deletor UserID
	Reason  *string
}

// Copy creates a deep copy
func (e *MessageDeleted) Copy() eventlog.Payload {
	cp := *e

	if e.Reason != nil {
		v := *e.Reason
		cp.Reason = &v
	}

	return &cp
}
