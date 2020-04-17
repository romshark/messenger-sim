package simulator_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/passhash"
	"github.com/romshark/messenger-sim/messenger/username"
	"github.com/romshark/messenger-sim/service/users"
	userssim "github.com/romshark/messenger-sim/service/users/simulator"

	"github.com/stretchr/testify/require"
)

type TestSetup struct {
	EventLog *eventlog.EventLog
	Service  users.Client
}

func NewTestSetup(t *testing.T) *TestSetup {
	l := eventlog.New()

	s, err := userssim.New(l, passhash.NoHash{})
	require.NoError(t, err)

	return &TestSetup{
		EventLog: l,
		Service:  s,
	}
}

func TestCreateNewUser(t *testing.T) {
	s := NewTestSetup(t)

	const (
		username    = username.Username("someusername")
		displayName = "someDisplayName"
	)

	newUser, err := s.Service.CreateNewUser(
		context.Background(),
		username,
		displayName,
		nil, // No avatar URL
		"testpassword",
	)
	require.NoError(t, err)
	require.WithinDuration(
		t,
		time.Now(),
		newUser.CreationTime,
		MaxCreationTimeDelta,
	)
	require.NotZero(t, newUser.ID)
	require.Equal(t, username, newUser.Username)
	require.Equal(t, displayName, newUser.DisplayName)
	require.Nil(t, newUser.AvatarURL)

	// Try to find the newly created user by ID
	foundUsers, err := s.Service.GetUsers(
		context.Background(),
		[]event.UserID{newUser.ID},
	)
	require.NoError(t, err)
	require.Len(t, foundUsers, 1)
	foundUser := foundUsers[0]
	require.Equal(t, newUser.CreationTime, foundUser.CreationTime)
	require.Equal(t, newUser.ID, foundUser.ID)
	require.Equal(t, username, foundUser.Username)
	require.Equal(t, displayName, foundUser.DisplayName)
	require.Nil(t, foundUser.AvatarURL)

	// Verify the event log
	buf := make([]eventlog.Event, 1)
	read, version, err := s.EventLog.Read(eventlog.Version(0), buf)
	require.NoError(t, err)
	require.Equal(t, eventlog.Version(1), version)
	require.Equal(t, 1, read)

	p := buf[0].Payload.(*event.UserCreated)
	require.Equal(t, newUser.ID, p.ID)
	require.Equal(t, username, p.Username)
	require.Equal(t, displayName, p.DisplayName)
	require.Nil(t, p.AvatarURL)
}

func TestCreateNewUser_UsernameReserved(t *testing.T) {
	s := NewTestSetup(t)

	const (
		username    = username.Username("someusername")
		displayName = "someDisplayName"
	)

	_, err := s.Service.CreateNewUser(
		context.Background(),
		username,
		displayName,
		nil,
		"testpassword1",
	)
	require.NoError(t, err)

	newUser, err := s.Service.CreateNewUser(
		context.Background(),
		username,
		displayName,
		nil,
		"testpassword2",
	)
	require.Error(t, err)
	require.True(t, errors.Is(err, users.ErrUsernameReserved))
	require.Nil(t, newUser)
}

const MaxCreationTimeDelta = time.Second
