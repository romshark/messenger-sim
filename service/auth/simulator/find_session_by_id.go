package simulator

import (
	"context"

	"github.com/romshark/messenger-sim/messenger/sessid"
	"github.com/romshark/messenger-sim/service/auth"
)

func (s *Simulator) FindSessionByID(
	ctx context.Context,
	sessionID sessid.SessionID,
) (*auth.Session, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	sess, ok := s.sessionsByID[sessionID]
	if !ok {
		return nil, nil
	}

	return &auth.Session{
		ID:           sessionID,
		User:         sess.user.id,
		IP:           sess.IP,
		UserAgent:    sess.UserAgent,
		CreationTime: sess.CreationTime,
	}, nil
}
