package auth

import (
	"context"
	"simulator/messenger/event"
	"simulator/messenger/sessid"
	"simulator/messenger/username"
	"time"
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
