package event

import "simulator/messenger/eventlog"

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
