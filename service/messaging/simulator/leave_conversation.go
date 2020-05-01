package simulator

import (
	"context"
	"fmt"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
)

func (s *Simulator) LeaveConversation(
	ctx context.Context,
	userID event.UserID,
	conversationID event.ConversationID,
) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, err := eventlog.TryPush(
		ctx,
		s.eventLog,
		func(retries int) (eventlog.Payload, error) {
			// Make sure the conversation exists
			var conv *conversation
			var ok bool
			if conv, ok = s.conversationsByID[conversationID]; !ok {
				return nil, fmt.Errorf(
					"conversation (%s) not found",
					conversationID.String(),
				)
			}

			// Make sure the user exists
			if _, ok := s.usersByID[userID]; !ok {
				return nil, fmt.Errorf(
					"user (%s) not found",
					userID.String(),
				)
			}

			// Make sure the user is part of the conversation
			if !func() bool {
				for _, participant := range conv.participants {
					if participant.id == userID {
						return true
					}
				}
				return false
			}() {
				return nil, fmt.Errorf(
					"user (%s) is not part of conversation (%s)",
					userID.String(),
					conversationID.String(),
				)
			}

			return &event.UserLeftConversation{
				User:         userID,
				Conversation: conversationID,
			}, nil
		},
		s.sync,
	)
	return err
}
