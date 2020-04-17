package simulator

import (
	"context"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/service/auth"
)

func (s *Simulator) ListSessionsForUser(
	ctx context.Context,
	userID event.UserID,
) ([]*auth.Session, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	u := func() *user {
		for _, u := range s.usersByUsername {
			if u.id == userID {
				return u
			}
		}
		return nil
	}()

	if u == nil {
		return nil, nil
	}

	l := make([]*auth.Session, 0, len(u.sessions))
	for _, s := range u.sessions {
		sess := s.Session
		l = append(l, &sess)
	}
	return l, nil
}
