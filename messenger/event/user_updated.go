package event

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/username"
)

type UserUpdated struct {
	User        UserID
	Username    *username.Username
	DisplayName *string
	AvatarURL   interface{}
}

// Copy creates a deep copy
func (e *UserUpdated) Copy() eventlog.Payload {
	cp := *e

	if e.Username != nil {
		v := *e.Username
		cp.Username = &v
	}

	if e.DisplayName != nil {
		v := *e.DisplayName
		cp.DisplayName = &v
	}

	switch v := e.AvatarURL.(type) {
	case *url.URL:
		if v != nil {
			u := *v
			cp.AvatarURL = &u
		}
	case nil:
	default:
		panic(fmt.Errorf("unexpected type: %s", reflect.TypeOf(v)))
	}

	return &cp
}
