package auth

import (
	"context"
	"time"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/sessid"
	"github.com/romshark/messenger-sim/messenger/username"
)

// Client represents the interface of an authentication service client
type Client interface {
	CreateSession(
		ctx context.Context,
		username username.Username,
		password string,
		ip string,
		userAgent string,
	) (*Session, error)

	DestroySession(
		ctx context.Context,
		sessionID sessid.SessionID,
	) error

	FindSessionByID(
		ctx context.Context,
		sessionID sessid.SessionID,
	) (*Session, error)

	ListSessionsForUser(
		ctx context.Context,
		userID event.UserID,
	) ([]*Session, error)
}

type Session struct {
	ID           sessid.SessionID
	User         event.UserID
	IP           string
	UserAgent    string
	CreationTime time.Time
}
