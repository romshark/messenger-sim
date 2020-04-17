package model

import (
	"simulator/messenger/event"
	"time"
)

type Session struct {
	UserID event.UserID

	ID           string    `json:"id"`
	IP           string    `json:"ip"`
	UserAgent    string    `json:"userAgent"`
	CreationTime time.Time `json:"creationTime"`
}
