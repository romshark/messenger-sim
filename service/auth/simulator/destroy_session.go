package simulator

import (
	"context"
	"fmt"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/sessid"
)

func (s *Simulator) DestroySession(
	ctx context.Context,
	sessionID sessid.SessionID,
) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, err := eventlog.TryPush(
		ctx,
		s.eventLog,
		s.projectionVersion,
		func(retries int) (eventlog.Payload, error) {
			// Make sure the session exists
			if _, ok := s.sessionsByID[sessionID]; !ok {
				return nil, fmt.Errorf("session (%s) not found", sessionID)
			}

			return &event.SessionDestroyed{
				Session: sessionID,
			}, nil
		},
		s.sync,
	)
	return err
}
