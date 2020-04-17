package event

import (
	"net/url"
	"simulator/messenger/eventlog"
	"simulator/messenger/username"
)

type UserCreated struct {
	ID           UserID
	Username     username.Username
	DisplayName  string
	AvatarURL    *url.URL
	PasswordHash string
}

// Copy creates a deep copy
func (e *UserCreated) Copy() eventlog.Payload {
	cp := *e

	if e.AvatarURL != nil {
		v := *e.AvatarURL
		cp.AvatarURL = &v
	}

	return &cp
}
