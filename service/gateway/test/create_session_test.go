package test

import (
	"context"
	"testing"

	"github.com/romshark/messenger-sim/messenger/username"
	"github.com/romshark/messenger-sim/service/gateway/client"
	"github.com/romshark/messenger-sim/service/gateway/client/model"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	s := NewSetup(t)
	require := require.New(t)

	const (
		pass = "testuser1password"
		name = "testuser"
	)

	u := s.CreateUser(
		t,
		username.Username(name),
		"Test User",
		nil,
		pass,
	)
	require.NoError(u.Init())

	const sessionsQuery = `query($uid:String!) {
		user(id:$uid) {
			sessions {
				id
				user { id }
				ip
				userAgent
				creationTime
			}
		}
	}`

	// Create first session
	s.Auth(t, u.Username, pass)

	var resp1 struct {
		User *model.User `json:"user"`
	}
	require.NoError(s.Query(
		context.Background(),
		&resp1,
		sessionsQuery,
		client.Args{"uid": u.ID},
	))

	require.Len(resp1.User.Sessions, 1)
	sess1 := resp1.User.Sessions[0]
	require.Equal(s.SessionID(), sess1.ID)
	require.Equal(s.UserID(), sess1.User.ID)
	require.NotZero(sess1.IP)
	require.NotZero(sess1.UserAgent)
	checkCreationTime(t, sess1.CreationTime)

	// Create second session
	s.Auth(t, u.Username, pass)

	var resp2 struct {
		User *model.User `json:"user"`
	}
	require.NoError(s.Query(
		context.Background(),
		&resp2,
		sessionsQuery,
		client.Args{"uid": u.ID},
	))

	require.Len(resp2.User.Sessions, 2)
	require.Equal(resp1.User.Sessions[0], resp2.User.Sessions[0])
	sess2 := resp2.User.Sessions[1]

	require.Equal(s.SessionID(), sess2.ID)
	require.NotEqual(sess1.ID, sess2.ID)

	require.Equal(s.UserID(), sess2.User.ID)
	require.NotZero(sess2.IP)
	require.NotZero(sess2.UserAgent)
	checkCreationTime(t, sess2.CreationTime)
}
