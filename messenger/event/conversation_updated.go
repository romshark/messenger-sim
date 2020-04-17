package event

import (
	"fmt"
	"net/url"
	"reflect"
	"simulator/messenger/eventlog"
)

type ConversationUpdated struct {
	Conversation ConversationID
	Editor       UserID
	Title        *string
	AvatarURL    interface{} // nil | *url.URL
}

// Copy creates a deep copy
func (e *ConversationUpdated) Copy() eventlog.Payload {
	cp := *e

	if e.Title != nil {
		v := *e.Title
		cp.Title = &v
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
