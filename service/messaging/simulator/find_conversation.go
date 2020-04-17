package simulator

import (
	"context"
	"simulator/messenger/event"
	"simulator/service/messaging"
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
