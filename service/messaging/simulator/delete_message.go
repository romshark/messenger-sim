package simulator

import (
	"context"
	"fmt"
	"simulator/messenger/event"
	"simulator/messenger/eventlog"
)

func (s *Simulator) DeleteMessage(
	ctx context.Context,
	messageID event.MessageID,
	deletorID event.UserID,
	reason *string,
) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, err := eventlog.TryPush(
		ctx,
		s.eventLog,
		s.projectionVersion,
		func(retries int) (eventlog.Payload, error) {
			// Make sure message exists
			if _, ok := s.messagesByID[messageID]; !ok {
				return nil, fmt.Errorf(
					"message (%s) not found",
					messageID.String(),
				)
			}

			// Make sure deletor exists
			if _, ok := s.usersByID[deletorID]; !ok {
				return nil, fmt.Errorf(
					"user (deletor) (%s) not found",
					deletorID.String(),
				)
			}

			e := &event.MessageDeleted{
				Message: messageID,
				Deletor: deletorID,
			}
			if reason != nil {
				v := *reason
				e.Reason = &v
			}
			return e, nil
		},
		s.sync,
	)
	return err
}
