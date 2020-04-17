package simulator

import (
	"context"
	"simulator/messenger/event"
	"simulator/service/messaging"
)

func (s *Simulator) GetMessageEdits(
	ctx context.Context,
	messageID event.MessageID,
) ([]*messaging.MessageEdit, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	m, ok := s.messagesByID[messageID]
	if !ok {
		return nil, nil
	}

	edits := make([]*messaging.MessageEdit, len(m.edits))
	for i, e := range m.edits {
		cp := *e
		edits[i] = &cp
	}
	return edits, nil
}
