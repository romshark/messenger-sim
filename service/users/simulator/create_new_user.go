package simulator

import (
	"context"
	"fmt"
	"net/url"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/id"
	"github.com/romshark/messenger-sim/messenger/username"
	"github.com/romshark/messenger-sim/service/users"
)

// CreateNewUser creates a new messenger user
func (s *Simulator) CreateNewUser(
	ctx context.Context,
	username username.Username,
	displayName string,
	avatarURL *url.URL,
	password string,
) (*users.User, error) {
	passwordHash, err := s.passwordHasher.Hash([]byte(password))
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	newID, err := id.New()
	if err != nil {
		return nil, fmt.Errorf("creating unique identifier: %w", err)
	}

	if err := username.Validate(); err != nil {
		return nil, fmt.Errorf("initializing username: %w", err)
	}

	newEvent := &event.UserCreated{
		ID:          event.UserID(newID),
		Username:    username,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	pushedEvent, err := eventlog.TryPush(
		ctx,
		s.eventLog,
		s.projectionVersion,
		func(retries int) (eventlog.Payload, error) {
			// Check invariants
			if _, ok := s.usersByUsername[username]; ok {
				return nil, users.ErrUsernameReserved
			}

			return newEvent, nil
		},
		s.sync,
	)
	if err != nil {
		return nil, err
	}
	return &users.User{
		ID:           event.UserID(newID),
		Username:     username,
		DisplayName:  displayName,
		AvatarURL:    avatarURL,
		CreationTime: pushedEvent.Time,
		PasswordHash: string(passwordHash),
	}, nil
}
