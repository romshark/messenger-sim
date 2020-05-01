package model

import (
	"github.com/romshark/messenger-sim/messenger/id"
	"github.com/romshark/messenger-sim/service/gateway/graph/model"
)

type Conversation struct {
	model.Conversation

	Participants []*User    `json:"participants"`
	Messages     []*Message `json:"messages"`
}

// Init initializes uninitialized fields
func (c *Conversation) Init() error {
	// Init ConversationID from ID
	if c.ID != "" && c.ConversationID.IsZero() {
		i, err := id.FromString(c.ID)
		if err != nil {
			return err
		}
		c.ConversationID = i
	}
	return nil
}
