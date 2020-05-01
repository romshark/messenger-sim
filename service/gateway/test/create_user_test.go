package test

import (
	"context"
	"testing"

	"github.com/romshark/messenger-sim/messenger/username"
	"github.com/romshark/messenger-sim/service/gateway/client"
	"github.com/romshark/messenger-sim/service/gateway/client/model"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	s := NewSetup(t)
	r := require.New(t)

	const (
		u1name   = "testuser1"
		u2name   = "testuser2"
		u1dpName = "Test User 1"
		u2dpName = "Test User 2"
		u1pass   = "testuser1password"
		u2pass   = "testuser2password"
	)

	u1 := s.CreateUser(
		t,
		username.Username(u1name),
		u1dpName,
		nil,
		u1pass,
	)
	u2 := s.CreateUser(
		t,
		username.Username(u2name),
		u2dpName,
		nil,
		u2pass,
	)

	checkUser := func(
		u *model.User,
		name username.Username,
		displayName string,
		password string,
	) {
		s.Auth(t, u.Username, password)

		var q struct{ User *model.User }
		r.NoError(s.Query(
			context.Background(),
			&q,
			`query($uid: String!) {
				user(id: $uid) {
					id
					username
					displayName
					creationTime
					avatarURL
					sessions { id }
					conversations { id }
				}
			}`,
			client.Args{"uid": u.ID},
		))

		x := q.User
		r.NotNil(x)
		r.Equal(u.ID, x.ID)
		r.Equal(name, username.Username(x.Username))
		r.Equal(displayName, x.DisplayName)
		if u.AvatarURL != nil {
			r.Equal(*u.AvatarURL, *x.AvatarURL)
		} else {
			r.Nil(x.AvatarURL)
		}
		r.Len(x.Sessions, 1)
		r.Equal(x.Sessions[0].ID, s.SessionID())
		r.Len(x.Conversations, 0)
	}

	checkUser(u1, u1name, u1dpName, u1pass)
	checkUser(u2, u2name, u2dpName, u2pass)
}
