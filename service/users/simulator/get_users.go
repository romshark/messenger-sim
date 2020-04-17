package simulator

import (
	"context"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/service/users"
)

// GetUsers returns user profiles
func (s *Simulator) GetUsers(
	ctx context.Context,
	ids []event.UserID,
) ([]*users.User, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	r := make([]*users.User, len(ids))
	for i := range r {
		if u, ok := s.usersByID[ids[i]]; ok {
			r[i] = u.Copy()
		}
	}

	return r, nil
}
