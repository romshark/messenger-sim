package simulator

import (
	"context"
	"fmt"
	"net/url"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/service/messaging"
)

func (s *Simulator) EditConversation(
	ctx context.Context,
	conversationID event.ConversationID,
	editorID event.UserID,
	title *string,
	avatarURL interface{},
) (*messaging.Conversation, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var conv *conversation

	_, err := eventlog.TryPush(
		ctx,
		s.eventLog,
		s.projectionVersion,
		func(retries int) (eventlog.Payload, error) {
			// Make sure the conversation exists
			var ok bool
			conv, ok = s.conversationsByID[conversationID]
			if !ok {
				return nil, fmt.Errorf(
					"conversation (%s) not found",
					conversationID.String(),
				)
			}

			// Make sure the editor exists
			if _, ok := s.usersByID[editorID]; !ok {
				return nil, fmt.Errorf(
					"user (editor) (%s) not found",
					editorID.String(),
				)
			}

			return &event.ConversationUpdated{
				Conversation: conversationID,
				Editor:       editorID,
				Title:        title,
				AvatarURL:    avatarURL,
			}, nil
		},
		s.sync,
	)
	if err != nil {
		return nil, err
	}

	newConv := &messaging.Conversation{
		ID:           conv.ID,
		CreationTime: conv.CreationTime,
	}

	if title != nil {
		newConv.Title = *title
	}
	if v, ok := avatarURL.(*url.URL); ok {
		newConv.AvatarURL = v
	}

	return newConv, nil
}
