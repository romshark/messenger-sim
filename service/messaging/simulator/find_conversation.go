package simulator

import (
	"context"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/service/messaging"
)

func (s *Simulator) FindConversation(
	ctx context.Context,
	id event.ConversationID,
) (*messaging.Conversation, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	c, ok := s.conversationsByID[id]
	if !ok {
		return nil, nil
	}
	return c.Conversation.Copy(), nil
}
