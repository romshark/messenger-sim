package simulator

import (
	"context"
	"fmt"
	"simulator/messenger/event"
	"simulator/service/messaging"
)

func (s *Simulator) GetMessages(
	ctx context.Context,
	conversationID event.ConversationID,
	afterID *event.MessageID,
	limit int,
) ([]*messaging.Message, error) {
	if limit < 1 {
		return nil, fmt.Errorf("invalid limit (%d)", limit)
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	c, ok := s.conversationsByID[conversationID]
	if !ok {
		return nil, nil
	}

	var slice []*message
	if afterID == nil {
		// Read at begin
		if limit > len(c.messages) {
			slice = c.messages
		} else {
			slice = c.messages[:limit]
		}
	} else {
		// Read after certain message
		indexOf := func() int {
			for i, m := range c.messages {
				if m.ID == *afterID {
					return i
				}
			}
			return -1
		}()

		if indexOf < 0 {
			return nil, fmt.Errorf(
				"message (%s) not found",
				afterID.String(),
			)
		}
		indexOf++
		if indexOf >= len(c.messages) {
			return nil, nil
		}

		tail := indexOf + limit
		if tail > len(c.messages) {
			tail = len(c.messages)
		}
		slice = c.messages[indexOf:tail]
	}

	messages := make([]*messaging.Message, len(slice))
	for i, msg := range slice {
		m := msg.Message
		messages[i] = &m
	}
	return messages, nil
}
