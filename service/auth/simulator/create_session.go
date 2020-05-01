package simulator

import (
	"context"
	"errors"
	"fmt"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/username"
	"github.com/romshark/messenger-sim/service/auth"
)

func (s *Simulator) CreateSession(
	ctx context.Context,
	username username.Username,
	password string,
	ip string,
	userAgent string,
) (*auth.Session, error) {
	newID, err := s.idGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generating id: %w", err)
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	var user *user

	pushedEvent, err := eventlog.TryPush(
		ctx,
		s.eventLog,
		func(retries int) (eventlog.Payload, error) {
			// Make sure the user exists
			var ok bool
			user, ok = s.usersByUsername[username]
			if !ok {
				return nil, ErrWrongCredentials
			}

			// Make sure the password matches
			ok, err := s.passwordComparer.Compare(
				[]byte(password),
				[]byte(user.passwordHash),
			)
			if err != nil {
				return nil, fmt.Errorf("comparing password: %w", err)
			}
			if !ok {
				return nil, ErrWrongCredentials
			}

			return &event.SessionCreated{
				ID:        newID,
				User:      user.id,
				IP:        ip,
				UserAgent: userAgent,
			}, nil
		},
		s.sync,
	)
	if err != nil {
		return nil, err
	}
	return &auth.Session{
		ID:           newID,
		User:         user.id,
		IP:           ip,
		UserAgent:    userAgent,
		CreationTime: pushedEvent.Time,
	}, nil
}

// ErrWrongCredentials indicates that either a wrong username
// or a wrong password has been passed
var ErrWrongCredentials = errors.New("wrong credentials")
