package model

import (
	"time"

	"github.com/romshark/messenger-sim/messenger/event"
)

type MessageEdit struct {
	EditorID event.UserID

	Time         time.Time `json:"time"`
	PreviousBody string    `json:"previousBody"`
}
