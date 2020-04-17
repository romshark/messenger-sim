package simulator

import (
	"context"
	"fmt"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/id"
	"github.com/romshark/messenger-sim/service/messaging"
)

func (s *Simulator) SendMessage(
	ctx context.Context,
	body string,
	conversationID event.ConversationID,
	senderID event.UserID,
) (*messaging.Message, error) {
	id, err := id.New()
	if err != nil {
		return nil, fmt.Errorf("creating identifier: %w", err)
	}
	newID := event.MessageID(id)

	// Validate body
	if body == "" {
		return nil, fmt.Errorf("invalid body (empty)")
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	pushedEvent, err := eventlog.TryPush(
		ctx,
		s.eventLog,
		s.projectionVersion,
		func(retries int) (eventlog.Payload, error) {
			// Make sure the conversation exists
			conv, ok := s.conversationsByID[conversationID]
			if !ok {
				return nil, fmt.Errorf(
					"conversation (%s) not found",
					conversationID.String(),
				)
			}

			// Make sure the sender exists
			if _, ok := s.usersByID[senderID]; !ok {
				return nil, fmt.Errorf(
					"user (sender) (%s) not found",
					senderID.String(),
				)
			}

			// Make sure the sender is part of the conversation
			if !func() bool {
				for _, participant := range conv.participants {
					if participant.id == senderID {
						return true
					}
				}
				return false
			}() {
				return nil, fmt.Errorf(
					"user (%s) isn't part of the conversation (%s)",
					senderID.String(),
					conversationID.String(),
				)
			}

			return &event.MessageSent{
				ID:           newID,
				Body:         body,
				Sender:       senderID,
				Conversation: conversationID,
			}, nil
		},
		s.sync,
	)
	if err != nil {
		return nil, err
	}

	return &messaging.Message{
		ID:           newID,
		Body:         body,
		Sender:       senderID,
		Conversation: conversationID,
		SendingTime:  pushedEvent.Time,
	}, nil
}
