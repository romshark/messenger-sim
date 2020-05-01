package simulator

import (
	"context"
	"fmt"
	"net/url"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/id"
	"github.com/romshark/messenger-sim/service/messaging"
)

func (s *Simulator) CreateConversation(
	ctx context.Context,
	title string,
	creatorID event.UserID,
	participants []event.UserID,
	avatarURL *url.URL,
) (*messaging.Conversation, error) {
	newID, err := id.New()
	if err != nil {
		return nil, fmt.Errorf("creating identifier: %w", err)
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	pushedEvent, err := eventlog.TryPush(
		ctx,
		s.eventLog,
		func(retries int) (eventlog.Payload, error) {
			// Make sure users exist
			for _, participantID := range participants {
				if _, ok := s.usersByID[participantID]; !ok {
					return nil, fmt.Errorf(
						"user (participant) (%s) not found",
						participantID.String(),
					)
				}
			}

			// Make sure the creator exists
			if _, ok := s.usersByID[creatorID]; !ok {
				return nil, fmt.Errorf(
					"user (creator) (%s) not found",
					creatorID.String(),
				)
			}

			return &event.ConversationCreated{
				ID:           event.ConversationID(newID),
				Creator:      creatorID,
				Participants: participants,
				Title:        title,
				AvatarURL:    avatarURL,
			}, nil
		},
		s.sync,
	)
	if err != nil {
		return nil, err
	}
	return &messaging.Conversation{
		ID:           event.ConversationID(newID),
		Title:        title,
		AvatarURL:    avatarURL,
		CreationTime: pushedEvent.Time,
	}, nil
}
