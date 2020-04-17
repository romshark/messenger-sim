package simulator

import (
	"fmt"
	"simulator/messenger/event"
	"simulator/messenger/eventlog"
	"simulator/messenger/passhash"
	"simulator/messenger/sessid"
	"simulator/messenger/username"
	"simulator/service/auth"
	"sync"
)

// Simulator represents an in-memory simulation of the authentication service
type Simulator struct {
	eventLog         *eventlog.EventLog
	idGenerator      *sessid.Generator
	passwordComparer passhash.PasswordComparer

	lock              sync.RWMutex
	projectionVersion eventlog.Version
	sessionsByID      map[sessid.SessionID]*session
	usersByUsername   map[username.Username]*user
}

// New creates a new instance of an authentication service simulator
func New(
	eventLog *eventlog.EventLog,
	idGenerator *sessid.Generator,
	passwordComparer passhash.PasswordComparer,
) (*Simulator, error) {
	if eventLog == nil {
		return nil, fmt.Errorf("missing event log")
	}
	if idGenerator == nil {
		return nil, fmt.Errorf("missing session identifier generator")
	}
	if passwordComparer == nil {
		return nil, fmt.Errorf("missing password comparer")
	}

	return &Simulator{
		eventLog:         eventLog,
		idGenerator:      idGenerator,
		passwordComparer: passwordComparer,

		sessionsByID:    make(map[sessid.SessionID]*session),
		usersByUsername: make(map[username.Username]*user),
	}, nil
}

type session struct {
	auth.Session
	user *user
}

type user struct {
	id           event.UserID
	passwordHash string
	sessions     map[sessid.SessionID]*session
}
