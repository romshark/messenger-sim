package model

import (
	"simulator/messenger/event"
	"time"
)

type MessageEdit struct {
	EditorID event.UserID

	Time         time.Time `json:"time"`
	PreviousBody string    `json:"previousBody"`
}
