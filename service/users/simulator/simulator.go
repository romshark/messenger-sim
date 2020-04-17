package simulator

import (
	"errors"
	"sync"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/passhash"
	"github.com/romshark/messenger-sim/messenger/username"
	"github.com/romshark/messenger-sim/service/users"
)

// Simulator represents an in-memory simulation of the users service
type Simulator struct {
	eventLog       *eventlog.EventLog
	passwordHasher passhash.PasswordHasher

	lock              sync.RWMutex
	projectionVersion eventlog.Version
	usersByUsername   map[username.Username]*users.User
	usersByID         map[event.UserID]*users.User
}

// New creates a new instance of a users service simulator
func New(
	eventLog *eventlog.EventLog,
	passwordHasher passhash.PasswordHasher,
) (*Simulator, error) {
	if eventLog == nil {
		return nil, errors.New("missing eventLog")
	}
	if passwordHasher == nil {
		return nil, errors.New("missing password hasher")
	}
	return &Simulator{
		eventLog:       eventLog,
		passwordHasher: passwordHasher,

		usersByUsername: make(map[username.Username]*users.User),
		usersByID:       make(map[event.UserID]*users.User),
	}, nil
}
