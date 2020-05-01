package simulator

import (
	"context"
	"fmt"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
)

func (s *Simulator) RemoveUserFromConversation(
	ctx context.Context,
	conversationID event.ConversationID,
	userID event.UserID,
	removerID event.UserID,
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

			// Make sure the removing user exists
			if _, ok := s.usersByID[removerID]; !ok {
				return nil, fmt.Errorf(
					"user (remover) (%s) not found",
					removerID.String(),
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
					"user (%s) is not part of the conversation (%s)",
					userID.String(),
					conversationID.String(),
				)
			}

			return &event.UserRemovedFromConversation{
				Remover:      removerID,
				Removed:      userID,
				Conversation: conversationID,
			}, nil
		},
		s.sync,
	)
	return err
}
