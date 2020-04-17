package eventlog

import (
	"errors"
	"fmt"
	"sync"

	"github.com/romshark/messenger-sim/messenger/id"
)

// Version represents a version identifier of the event log
// which is derived from the log's length
type Version uint64

// EventLog represents an in-memory event log implementation
type EventLog struct {
	lock sync.RWMutex
	log  []Event
	subs map[id.ID]chan<- Version
}

// New creates a new event log instance
func New() *EventLog {
	return &EventLog{
		subs: make(map[id.ID]chan<- Version),
	}
}

// Version returns the current version of the log
func (s *EventLog) Version() Version {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.version()
}

// Push pushes a new event onto the log
func (s *EventLog) Push(payload Payload) (
	newVersion Version,
	pushedEvent Event,
	err error,
) {
	if payload == nil {
		newVersion = s.Version()
		return
	}

	e, err := NewEvent(payload)
	if err != nil {
		err = fmt.Errorf("creating a new event instance: %w", err)
		return
	}

	s.lock.Lock()
	s.log = append(s.log, e)
	newVersion = s.version()
	s.lock.Unlock()

	// Broadcast change
	go func() {
		s.lock.RLock()
		defer s.lock.RUnlock()
		s.broadcast()
	}()

	pushedEvent = e
	pushedEvent.Payload = payload // avoid returning the payload copy

	return
}

// ErrMismatchingVersion indicates that the assumed version
// doesn't match up with the current version of the event log
var ErrMismatchingVersion = errors.New("mismatching version")

// CheckPush ensure that the assumed version is the same
// as the current version before pushing a new event onto the log
func (s *EventLog) CheckPush(
	assumedVersion Version,
	payload Payload,
) (
	newVersion Version,
	pushedEvent Event,
	err error,
) {
	if payload == nil {
		newVersion = s.Version()
		return
	}

	e, err := NewEvent(payload)
	if err != nil {
		err = fmt.Errorf("creating a new event instance: %w", err)
		return
	}

	s.lock.Lock()
	if s.version() != assumedVersion {
		s.lock.Unlock()
		err = ErrMismatchingVersion
		return
	}
	s.log = append(s.log, e)
	newVersion = s.version()
	s.lock.Unlock()

	// Broadcast change
	go func() {
		s.lock.RLock()
		defer s.lock.RUnlock()
		s.broadcast()
	}()

	pushedEvent = e
	pushedEvent.Payload = payload // avoid returning the payload copy

	return
}

// Read tries to fill the given buffer with events pushed at the given version
func (s *EventLog) Read(
	after Version,
	buf []Event,
) (
	read int,
	logVersion Version,
	err error,
) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	v := Version(len(s.log))
	switch {
	case after == v:
		return 0, v, nil
	case after > v:
		return 0, 0, fmt.Errorf(
			"version (%d) ahead of actual log version (%d)",
			after, v,
		)
	}

	var section []Event
	tail := int(after) + len(buf)
	if tail > len(s.log) {
		tail = len(s.log)
	}
	section = s.log[after:tail]
	for i, e := range section {
		buf[i] = Event{
			ID:      e.ID,
			Time:    e.Time,
			Payload: e.Payload.Copy(),
		}
	}
	return len(section), s.version(), nil
}

// NoLimit is used for EventLog.Read to indicate unlimited reading capacity
const NoLimit = uint32(0)

// Subscribe creates a new subscription pushing
// newly pushed events onto the channel in a non-blocking fashion
func (s *EventLog) Subscribe() (Subscription, error) {
	c := make(chan Version)
	s.lock.Lock()
	defer s.lock.Unlock()
	subID, err := id.New()
	if err != nil {
		return Subscription{}, fmt.Errorf("creating subscription id: %w", err)
	}
	s.subs[subID] = c

	cancel := func() {
		s.lock.Lock()
		defer s.lock.Unlock()
		close(s.subs[subID])
		delete(s.subs, subID)
	}

	return Subscription{cancel, c}, nil
}

func (s *EventLog) version() Version { return Version(len(s.log)) }

func (s *EventLog) broadcast() (sent int) {
	v := s.version()
	for _, sub := range s.subs {
		select {
		case sub <- v:
		default:
		}
		sent++
	}
	return
}

// Subscription represents a subscription to the event log
type Subscription struct {
	cancel func()
	c      chan Version
}

// C returns the update notification channel
func (s Subscription) C() <-chan Version { return s.c }

// Cancel cancels the subscription closing the channel
func (s Subscription) Cancel() { s.cancel() }
