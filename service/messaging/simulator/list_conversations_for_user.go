package simulator

import (
	"context"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/service/messaging"
)

func (s *Simulator) ListConversationsForUser(
	ctx context.Context,
	userID event.UserID,
) ([]*messaging.Conversation, error) {
	u, ok := s.usersByID[userID]
	if !ok {
		return nil, nil
	}

	l := make([]*messaging.Conversation, 0, len(u.joinedConversations))
	for _, rel := range u.joinedConversations {
		l = append(l, rel.conversation.Copy())
	}
	return l, nil
}
