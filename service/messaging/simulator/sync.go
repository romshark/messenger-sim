package simulator

import (
	"fmt"
	"net/url"
	"reflect"
	"simulator/messenger/event"
	"simulator/messenger/eventlog"
	"simulator/service/messaging"
)

// sync synchronizes the service against the event log
func (s *Simulator) sync() (eventlog.Version, error) {
	latestVersion, err := eventlog.Scan(
		s.eventLog,
		s.projectionVersion,
		nil,
		func(e eventlog.Event) bool {
			switch p := e.Payload.(type) {
			case *event.UserCreated:
				s.applyUserCreated(e, p)
			case *event.MessageSent:
				s.applyMessageSent(e, p)
			case *event.MessageDeleted:
				s.applyMessageDeleted(e, p)
			case *event.MessageEdited:
				s.applyMessageEdited(e, p)
			case *event.MessageRead:
				s.applyMessageRead(e, p)
			case *event.ConversationCreated:
				s.applyConversationCreated(e, p)
			case *event.ConversationUpdated:
				s.applyConversationUpdated(e, p)
			case *event.UserJoinedConversation:
				s.applyUserJoinedConversation(e, p)
			case *event.UserLeftConversation:
				s.applyUserLeftConversation(e, p)
			case *event.UserRemovedFromConversation:
				s.applyUserRemovedFromConversation(e, p)
			}
			return true
		},
	)
	s.projectionVersion = latestVersion
	return latestVersion, err
}

func (s *Simulator) applyUserCreated(e eventlog.Event, p *event.UserCreated) {
	s.usersByID[p.ID] = &user{
		joinedConversations: make(map[event.ConversationID]relUserConv),
	}
}

func (s *Simulator) applyMessageSent(e eventlog.Event, p *event.MessageSent) {
	c := s.conversationsByID[p.Conversation]
	newMessage := &message{
		Message: messaging.Message{
			ID:           p.ID,
			Body:         p.Body,
			Sender:       p.Sender,
			SendingTime:  e.Time,
			Conversation: p.Conversation,
		},
		conversation: c,
	}
	c.messages = append(c.messages, newMessage)
	s.messagesByID[p.ID] = newMessage
}

func (s *Simulator) applyMessageDeleted(e eventlog.Event, p *event.MessageDeleted) {
	m := s.messagesByID[p.Message]
	delete(s.messagesByID, p.Message)
	for i, msg := range m.conversation.messages {
		if msg.ID == p.Message {
			m.conversation.messages = append(
				m.conversation.messages[:i],
				m.conversation.messages[i+1:]...,
			)
		}
	}
}

func (s *Simulator) applyMessageEdited(e eventlog.Event, p *event.MessageEdited) {
	m := s.messagesByID[p.Message]
	m.edits = append(m.edits, &messaging.MessageEdit{
		Editor:       p.Editor,
		PreviousBody: m.Body,
		Time:         e.Time,
	})
	m.Body = p.Body
}

func (s *Simulator) applyMessageRead(
	e eventlog.Event,
	p *event.MessageRead,
) {
	panic("not yet implemented")
}

func (s *Simulator) applyConversationCreated(
	e eventlog.Event,
	p *event.ConversationCreated,
) {
	newConv := &conversation{
		Conversation: messaging.Conversation{
			ID:           p.ID,
			Title:        p.Title,
			AvatarURL:    p.AvatarURL,
			CreationTime: e.Time,
		},
	}
	s.conversationsByID[p.ID] = newConv
	for _, participantID := range p.Participants {
		s.usersByID[participantID].joinedConversations[p.ID] = relUserConv{
			joined:       e.Time,
			conversation: newConv,
		}
	}
}

func (s *Simulator) applyConversationUpdated(
	e eventlog.Event,
	p *event.ConversationUpdated,
) {
	c := s.conversationsByID[p.Conversation]

	// Update title if necessary
	if p.Title != nil {
		c.Title = *p.Title
	}

	// Update avatar URL if necessary
	switch v := p.AvatarURL.(type) {
	case *url.URL:
		if v == nil {
			// Erase
			c.AvatarURL = nil
		} else {
			// Change
			u := *v
			c.AvatarURL = &u
		}
	case nil:
		// Unchanged
	default:
		panic(fmt.Errorf("unexpected type: %s", reflect.TypeOf(v)))
	}
}

func (s *Simulator) applyUserJoinedConversation(
	e eventlog.Event,
	p *event.UserJoinedConversation,
) {
	c := s.conversationsByID[p.Conversation]
	u := s.usersByID[p.User]
	c.participants[u.id] = u
	u.joinedConversations[c.ID] = relUserConv{
		joined:       e.Time,
		conversation: c,
	}
}

func (s *Simulator) applyUserLeftConversation(
	e eventlog.Event,
	p *event.UserLeftConversation,
) {
	c := s.conversationsByID[p.Conversation]
	u := s.usersByID[p.User]
	delete(c.participants, u.id)
	delete(u.joinedConversations, c.ID)
}

func (s *Simulator) applyUserRemovedFromConversation(
	e eventlog.Event,
	p *event.UserRemovedFromConversation,
) {
	c := s.conversationsByID[p.Conversation]
	u := s.usersByID[p.Removed]
	delete(c.participants, u.id)
	delete(u.joinedConversations, c.ID)
}
