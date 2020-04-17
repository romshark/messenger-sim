package model

import (
	"simulator/messenger/event"
	"time"
)

type User struct {
	UserID event.UserID

	ID           string    `json:"id"`
	Username     string    `json:"username"`
	DisplayName  string    `json:"displayName"`
	CreationTime time.Time `json:"creationTime"`
	AvatarURL    *string   `json:"avatarURL"`
}
