package test

import (
	"context"
	"testing"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/username"
	"github.com/romshark/messenger-sim/service/gateway/client"
	"github.com/romshark/messenger-sim/service/gateway/client/model"
	"github.com/stretchr/testify/require"
)

func TestCreateConversation(t *testing.T) {
	s := NewSetup(t)
	r := require.New(t)

	const (
		u1pass = "testuser1password"
		u2pass = "testuser2password"
	)

	u1 := s.CreateUser(
		t,
		username.Username("testuser1"),
		"Test User 1",
		nil,
		u1pass,
	)
	r.NoError(u1.Init())

	u2 := s.CreateUser(
		t,
		username.Username("testuser2"),
		"Test User 2",
		nil,
		u2pass,
	)
	r.NoError(u2.Init())

	s.Auth(t, u1.Username, u1pass)
	c1 := s.CreateConversation(
		t,
		"testconv",
		[]event.UserID{u1.UserID, u2.UserID},
		nil,
	)

	checkConversation := func(
		userName string,
		userPassword string,
		userID event.UserID,
		conv *model.Conversation,
	) {
		s.Auth(t, userName, userPassword)

		var q struct{ User *model.User }
		r.NoError(s.Query(
			context.Background(),
			&q,
			`query($uid: String!) {
				user(id: $uid) {
					conversations {
						id
						title
						avatarURL
						creationTime
						messages(limit: 10) { id }
						participants { id }
					}
				}
			}`,
			client.Args{"uid": userID},
		))

		r.NotNil(q.User)
		r.Len(q.User.Conversations, 1)

		x := q.User.Conversations[0]

		r.NotNil(x)
		r.Equal(conv.ID, x.ID)
		r.Equal(conv.Title, x.Title)
		if conv.AvatarURL != nil {
			r.Equal(*conv.AvatarURL, *x.AvatarURL)
		} else {
			r.Nil(x.AvatarURL)
		}
		r.Len(x.Messages, 0)
		r.Len(x.Participants, len(conv.Participants))
		for _, p := range conv.Participants {
			r.True(func(expected *model.User) bool {
				for _, p := range x.Participants {
					if p.ID == expected.ID {
						return true
					}
				}
				return false
			}(p))
		}
	}

	checkConversation(u1.Username, u1pass, u1.UserID, c1)
	checkConversation(u2.Username, u2pass, u2.UserID, c1)
}
