package users

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/username"
)

// Client represents the interface of a users service client
type Client interface {
	CreateNewUser(
		ctx context.Context,
		username username.Username,
		displayName string,
		avatarURL *url.URL,
		password string,
	) (*User, error)

	GetUsers(
		ctx context.Context,
		ids []event.UserID,
	) ([]*User, error)
}

// User represents a user profile aggregate state
type User struct {
	ID           event.UserID
	Username     username.Username
	DisplayName  string
	CreationTime time.Time
	AvatarURL    *url.URL
	PasswordHash string
}

// Copy returns a deep copy of the user object
func (u *User) Copy() *User {
	cp := *u

	if cp.AvatarURL != nil {
		v := *cp.AvatarURL
		cp.AvatarURL = &v
	}

	return &cp
}

// ErrUsernameReserved indicates that the given username
// has already been reserved by another user
var ErrUsernameReserved = errors.New("username is already reserved")
