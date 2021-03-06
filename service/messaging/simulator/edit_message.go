package simulator

import (
	"context"
	"fmt"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/service/messaging"
)

func (s *Simulator) EditMessage(
	ctx context.Context,
	messageID event.MessageID,
	editorID event.UserID,
	body string,
) (*messaging.Message, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var msg *message

	_, err := eventlog.TryPush(
		ctx,
		s.eventLog,
		func(retries int) (eventlog.Payload, error) {
			// Make sure the message exists
			var ok bool
			msg, ok = s.messagesByID[messageID]
			if !ok {
				return nil, fmt.Errorf(
					"message (%s) not found",
					messageID.String(),
				)
			}

			// Make sure the body changed
			if msg.Body == body {
				return nil, fmt.Errorf("message body unchanged")
			}

			// Make sure the editor exists
			if _, ok := s.usersByID[editorID]; !ok {
				return nil, fmt.Errorf(
					"user (editor) (%s) not found",
					editorID.String(),
				)
			}

			return &event.MessageEdited{
				Message: messageID,
				Editor:  editorID,
				Body:    body,
			}, nil
		},
		s.sync,
	)
	if err != nil {
		return nil, err
	}
	return &messaging.Message{
		ID:           msg.ID,
		Body:         body,
		Sender:       msg.Sender,
		Conversation: msg.conversation.ID,
		SendingTime:  msg.SendingTime,
	}, nil
}
