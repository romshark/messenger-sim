package simulator

import (
	"errors"
	"sync"
	"time"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/service/messaging"
)

// Simulator represents an in-memory simulation of the messaging service
type Simulator struct {
	eventLog *eventlog.EventLog

	lock              sync.RWMutex
	projectionVersion eventlog.Version
	conversationsByID map[event.ConversationID]*conversation
	messagesByID      map[event.MessageID]*message
	usersByID         map[event.UserID]*user
}

// New creates a new instance of a messaging service simulator
func New(eventLog *eventlog.EventLog) (*Simulator, error) {
	if eventLog == nil {
		return nil, errors.New("missing eventLog")
	}
	return &Simulator{
		eventLog:          eventLog,
		conversationsByID: make(map[event.ConversationID]*conversation),
		messagesByID:      make(map[event.MessageID]*message),
		usersByID:         make(map[event.UserID]*user),
	}, nil
}

type conversation struct {
	messaging.Conversation
	messages     []*message
	participants map[event.UserID]*user
}

type message struct {
	messaging.Message
	edits        []*messaging.MessageEdit
	conversation *conversation
}

type user struct {
	id                  event.UserID
	joinedConversations map[event.ConversationID]relUserConv
}

// relUserConv defines a user-conversation relationship
type relUserConv struct {
	joined       time.Time
	conversation *conversation
}
