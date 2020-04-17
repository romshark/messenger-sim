package eventlog_test

import (
	"context"
	"testing"
	"time"

	"github.com/romshark/messenger-sim/messenger/eventlog"

	"github.com/stretchr/testify/require"
)

func TestTryPush(t *testing.T) {
	l := eventlog.New()

	push(t, l, "first", V(1), noAssume, noError)

	newEvent := &testEvent{Name: "new"}

	txCalls := 0
	var lastRetries int
	tx := func(retries int) (eventlog.Payload, error) {
		txCalls++
		lastRetries = retries
		return newEvent, nil
	}

	syncCalls := 0
	sync := func() (eventlog.Version, error) {
		syncCalls++
		return V(1), nil
	}

	pushedEvent, err := eventlog.TryPush(
		context.Background(),
		l,
		V(0), // Assume outdated version
		tx,
		sync,
	)
	require.NoError(t, err)
	require.NotZero(t, pushedEvent.ID)
	require.WithinDuration(
		t,
		time.Now(),
		pushedEvent.Time,
		time.Second,
	)
	require.Equal(t, 2, txCalls)
	require.Equal(t, 2, syncCalls, "expect one retry sync call and one final")
	require.Equal(t, 1, lastRetries)
}
