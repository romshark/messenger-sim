package simulator

import (
	"context"

	"github.com/romshark/messenger-sim/messenger/event"
)

func (s *Simulator) ListParticipants(
	ctx context.Context,
	conversationID event.ConversationID,
) ([]event.UserID, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	conv, ok := s.conversationsByID[conversationID]
	if !ok {
		return nil, nil
	}

	l := make([]event.UserID, 0, len(conv.participants))
	for i := range conv.participants {
		l = append(l, i)
	}
	return l, nil
}
