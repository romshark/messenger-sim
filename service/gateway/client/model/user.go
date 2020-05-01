package model

import (
	"github.com/romshark/messenger-sim/messenger/id"
	"github.com/romshark/messenger-sim/service/gateway/graph/model"
)

type User struct {
	model.User

	Sessions      []*Session      `json:"sessions"`
	Conversations []*Conversation `json:"conversations"`
}

// Init initializes uninitialized fields
func (u *User) Init() error {
	// Init UserID from ID
	if u.ID != "" && u.UserID.IsZero() {
		i, err := id.FromString(u.ID)
		if err != nil {
			return err
		}
		u.UserID = i
	}
	return nil
}
