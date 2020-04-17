package messaging

import (
	"context"
	"net/url"
	"time"

	"github.com/romshark/messenger-sim/messenger/event"
)

// Client represents the interface of a messaging service client
type Client interface {
	CreateConversation(
		ctx context.Context,
		title string,
		creatorID event.UserID,
		participants []event.UserID,
		avatarURL *url.URL,
	) (*Conversation, error)

	DeleteMessage(
		ctx context.Context,
		messageID event.MessageID,
		deletorID event.UserID,
		reason *string,
	) error

	EditConversation(
		ctx context.Context,
		conversationID event.ConversationID,
		editorID event.UserID,
		title *string,
		avatarURL interface{},
	) (*Conversation, error)

	EditMessage(
		ctx context.Context,
		messageID event.MessageID,
		editorID event.UserID,
		body string,
	) (*Message, error)

	FindConversation(
		ctx context.Context,
		id event.ConversationID,
	) (*Conversation, error)

	GetMessageEdits(
		ctx context.Context,
		messageID event.MessageID,
	) ([]*MessageEdit, error)

	GetMessages(
		ctx context.Context,
		conversationID event.ConversationID,
		afterID *event.MessageID,
		limit int,
	) ([]*Message, error)

	LeaveConversation(
		ctx context.Context,
		userID event.UserID,
		conversationID event.ConversationID,
	) error

	ListConversationsForUser(
		ctx context.Context,
		userID event.UserID,
	) ([]*Conversation, error)

	ListParticipants(
		ctx context.Context,
		conversationID event.ConversationID,
	) ([]event.UserID, error)

	RemoveUserFromConversation(
		ctx context.Context,
		conversationID event.ConversationID,
		userID event.UserID,
		removerID event.UserID,
	) error

	SendMessage(
		ctx context.Context,
		body string,
		conversationID event.ConversationID,
		senderID event.UserID,
	) (*Message, error)
}

type Conversation struct {
	ID           event.ConversationID
	Title        string
	AvatarURL    *url.URL
	CreationTime time.Time
}

func (c *Conversation) Copy() *Conversation {
	cp := *c

	if c.AvatarURL != nil {
		v := *c.AvatarURL
		cp.AvatarURL = &v
	}

	return &cp
}

type Message struct {
	ID           event.MessageID
	Body         string
	Sender       event.UserID
	Conversation event.ConversationID
	SendingTime  time.Time
}

type MessageEdit struct {
	Editor       event.UserID
	Time         time.Time
	PreviousBody string
}
