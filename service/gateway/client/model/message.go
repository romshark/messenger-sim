package model

import (
	"github.com/romshark/messenger-sim/messenger/id"
	"github.com/romshark/messenger-sim/service/gateway/graph/model"
)

type Message struct {
	model.Message

	Sender       []*User         `json:"sender"`
	Conversation []*Conversation `json:"conversation"`
}

// Init initializes uninitialized fields
func (m *Message) Init() error {
	// Init MessageID from ID
	if m.ID != "" && m.MessageID.IsZero() {
		i, err := id.FromString(m.ID)
		if err != nil {
			return err
		}
		m.MessageID = i
	}
	return nil
}
