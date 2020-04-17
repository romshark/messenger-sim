package eventlog

import (
	"errors"
	"time"

	"github.com/romshark/messenger-sim/messenger/id"
)

// EventID represents a unique event identifier
type EventID id.ID

// Event defines the basic information about an event
// such as its unique identifier and creation time
type Event struct {
	ID      EventID
	Time    time.Time
	Payload Payload
}

// NewEvent creates a new event with a copy of the given payload
func NewEvent(payload Payload) (Event, error) {
	if payload == nil {
		return Event{}, errors.New("missing payload")
	}
	newID, err := id.New()
	if err != nil {
		return Event{}, err
	}
	return Event{
		ID:      EventID(newID),
		Time:    time.Now().UTC(),
		Payload: payload.Copy(),
	}, nil
}

// Payload defines the interface of an event's payload
type Payload interface {
	Copy() Payload
}
