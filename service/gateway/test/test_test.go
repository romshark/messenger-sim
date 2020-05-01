package test

import (
	"context"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/passhash"
	"github.com/romshark/messenger-sim/messenger/sessid"
	uname "github.com/romshark/messenger-sim/messenger/username"
	authsim "github.com/romshark/messenger-sim/service/auth/simulator"
	"github.com/romshark/messenger-sim/service/gateway"
	"github.com/romshark/messenger-sim/service/gateway/client"
	"github.com/romshark/messenger-sim/service/gateway/client/model"
	messagingsim "github.com/romshark/messenger-sim/service/messaging/simulator"
	userssim "github.com/romshark/messenger-sim/service/users/simulator"

	"github.com/stretchr/testify/require"
)

type Setup struct {
	*client.Client
	GatewayServer *gateway.Server
}

func NewSetup(t *testing.T) *Setup {
	l := eventlog.New()

	passHashComparer := passhash.NoHash{}

	usersService, err := userssim.New(l, passHashComparer)
	require.NoError(t, err)

	sessIDGen, err := sessid.NewGenerator(128)
	require.NoError(t, err)

	authService, err := authsim.New(l, sessIDGen, passHashComparer)
	require.NoError(t, err)

	messagingService, err := messagingsim.New(l)
	require.NoError(t, err)

	gatewayServer, err := gateway.NewServer(
		usersService,
		authService,
		messagingService,
	)
	require.NoError(t, err)

	s := httptest.NewTLSServer(gatewayServer)
	httpClient := s.Client()

	j, err := cookiejar.New(nil)
	require.NoError(t, err)
	httpClient.Jar = j

	return &Setup{
		GatewayServer: gatewayServer,
		Client:        client.NewClient(httpClient, s.URL+"/g"),
	}
}

func (s *Setup) Auth(
	t *testing.T,
	username string,
	password string,
) {
	require.NoError(t, s.Client.Auth(
		context.Background(),
		username,
		password,
	))
	require.NotZero(t, s.Client.SessionID())
	require.NotZero(t, s.Client.UserID())
}

func (s *Setup) CreateUser(
	t *testing.T,
	username uname.Username,
	displayName string,
	avatarURL *url.URL,
	password string,
) *model.User {
	r := require.New(t)

	x, err := s.Client.CreateUser(
		context.Background(),
		username,
		displayName,
		avatarURL,
		password,
		`id
		username
		displayName
		creationTime
		avatarURL
		sessions { id }
		conversations { id }`,
	)
	r.NoError(err)

	r.NotNil(x)
	r.NotZero(x.ID)
	r.Equal(username, uname.Username(x.Username))
	r.Equal(displayName, x.DisplayName)
	checkCreationTime(t, x.CreationTime)
	checkAvatarURL(t, avatarURL, x.AvatarURL)
	r.Len(x.Sessions, 0)
	r.Len(x.Conversations, 0)

	return x
}

func (s *Setup) CreateConversation(
	t *testing.T,
	title string,
	participants []event.UserID,
	avatarURL *url.URL,
) *model.Conversation {
	r := require.New(t)

	x, err := s.Client.CreateConversation(
		context.Background(),
		title,
		participants,
		avatarURL,
		`id
		title
		avatarURL
		participants { id }
		messages(limit: 1) { id }
		creationTime`,
	)
	r.NoError(err)

	r.NotNil(x)
	r.NotZero(x.ID)
	r.Equal(title, x.Title)
	checkAvatarURL(t, avatarURL, x.AvatarURL)
	checkCreationTime(t, x.CreationTime)
	r.Len(x.Messages, 0)

	r.Len(x.Participants, len(participants))

	expectedPartIDs := make([]string, len(participants))
	for i, id := range participants {
		expectedPartIDs[i] = id.String()
	}

	for _, p := range x.Participants {
		r.Contains(expectedPartIDs, p.ID)
	}

	return x
}

func checkCreationTime(t *testing.T, actual time.Time) {
	require.WithinDuration(t, time.Now(), actual, time.Second)
}

func checkAvatarURL(t *testing.T, expected *url.URL, actual *string) {
	if expected != nil {
		require.Equal(t, expected.String(), *actual)
	} else {
		require.Nil(t, actual)
	}
}
