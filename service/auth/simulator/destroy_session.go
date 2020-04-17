package simulator

import (
	"context"
	"fmt"
	"simulator/messenger/event"
	"simulator/messenger/eventlog"
	"simulator/messenger/sessid"
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
