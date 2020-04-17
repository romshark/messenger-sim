package eventlog_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"simulator/messenger/eventlog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestPush tests Push
func TestPush(t *testing.T) {
	l := eventlog.New()
	require.Equal(t, V(0), l.Version())

	push(t, l, "first", V(1), noAssume, noError)
	push(t, l, "second", V(2), noAssume, noError)
}

// TestCheckPush tests CheckPush
func TestCheckPush(t *testing.T) {
	l := eventlog.New()

	push(t, l, "first", V(1), V(0), noError)
	push(t, l, "second", V(2), V(1), noError)

	v := V(2)

	// Version above current
	push(t, l, "potential third", v, v+1, eventlog.ErrMismatchingVersion)

	// Version below current
	push(t, l, "potential third", v, v-1, eventlog.ErrMismatchingVersion)
}

// TestPushNil makes sure Push and CheckPush won't append a nil payload
func TestPushNil(t *testing.T) {
	l := eventlog.New()
	require.Equal(t, V(0), l.Version())

	v, e, err := l.Push(nil)
	require.NoError(t, err)
	require.Equal(t, V(0), v)
	require.Zero(t, e)

	v, e, err = l.CheckPush(V(0), nil)
	require.NoError(t, err)
	require.Equal(t, V(0), v)
	require.Zero(t, e)
}

// TestRead tests Read
func TestRead(t *testing.T) {
	l := eventlog.New()

	expected := []string{"first", "second", "third", "fourth", "fifth"}

	for i, name := range expected {
		push(t, l, name, V(i+1), V(i), noError)
	}

	// Read 0-2
	expectedSlice := expected[:3]
	buf := make([]eventlog.Event, 3)

	read, version, err := l.Read(V(0), buf)
	require.NoError(t, err)
	require.Equal(t, V(len(expected)), version)
	require.Equal(t, 3, read)

	idRegister := checkExpected(t, expectedSlice, buf[:read], nil)

	// Read 3-4
	expectedSlice = expected[3:]

	read, version, err = l.Read(V(3), buf)
	require.NoError(t, err)
	require.Equal(t, V(len(expected)), version)
	require.Equal(t, 2, read)

	checkExpected(t, expectedSlice, buf[:read], idRegister)
}

// TestReadNilBuffer makes sure Read functions properly
// when the given buffer is nil
func TestReadNilBuffer(t *testing.T) {
	l := eventlog.New()

	// Read 0-2
	push(t, l, "first", V(1), V(0), noError)
	push(t, l, "second", V(2), V(1), noError)
	push(t, l, "third", V(3), V(2), noError)

	read, version, err := l.Read(V(0), nil)
	require.NoError(t, err)
	require.Zero(t, read)
	require.Equal(t, V(3), version)

	zeroBuf := make([]eventlog.Event, 0)

	read, version, err = l.Read(V(0), zeroBuf)
	require.NoError(t, err)
	require.Zero(t, read)
	require.Equal(t, V(3), version)
}

// TestReadAtLastVersion makes sure Read successfuly reads 0
// events when reading at the latest log version
func TestReadAtLastVersion(t *testing.T) {
	l := eventlog.New()

	// Read 0-2
	push(t, l, "first", V(1), V(0), noError)
	push(t, l, "second", V(2), V(1), noError)
	push(t, l, "third", V(3), V(2), noError)

	buf := make([]eventlog.Event, 10)

	read, version, err := l.Read(V(3), buf)
	require.NoError(t, err)
	require.Zero(t, read)
	require.Equal(t, V(3), version)
	for _, e := range buf {
		require.Zero(t, e.ID)
		require.Zero(t, e.Time)
		require.Nil(t, e.Payload)
	}
}

// TestSubscribe tests Subscribe
func TestSubscribe(t *testing.T) {
	l := eventlog.New()

	s, err := l.Subscribe()
	defer s.Cancel()
	require.NoError(t, err)

	const awaitTimeout = time.Second

	awaitNotification := func(expectedVersion V) <-chan error {
		c := make(chan error, 1)
		go func(c chan<- error) {
			start := time.Now()

			// Timeout when waiting for too long
			ctx, cancel := context.WithTimeout(context.Background(), awaitTimeout)
			defer cancel()

			select {
			case v := <-s.C():
				if expectedVersion != v {
					c <- fmt.Errorf(
						"unexpected version: %d (expected: %d)",
						v, expectedVersion,
					)
					return
				}
				c <- nil
			case <-ctx.Done():
				c <- fmt.Errorf(
					"timed out after %s, missing notification",
					time.Since(start),
				)
			}
		}(c)
		return c
	}

	c := awaitNotification(V(1))
	time.Sleep(10 * time.Millisecond)
	push(t, l, "first", V(1), V(0), noError)
	require.NoError(t, <-c)

	c = awaitNotification(V(2))
	time.Sleep(10 * time.Millisecond)
	push(t, l, "second", V(2), V(1), noError)
	require.NoError(t, <-c)

	c = awaitNotification(V(3))
	time.Sleep(10 * time.Millisecond)
	push(t, l, "third", V(3), V(2), noError)
	require.NoError(t, <-c)

	push(t, l, "fourth", V(4), V(3), noError)
	time.Sleep(10 * time.Millisecond)

	// Expect the skipped notification to be lost
	c = awaitNotification(V(5))
	time.Sleep(10 * time.Millisecond)
	push(t, l, "fifth", V(5), V(4), noError)
	require.NoError(t, <-c)
}

// TestSubscribeCancel tests subscription cancelation
func TestSubscribeCancel(t *testing.T) {
	l := eventlog.New()

	s, err := l.Subscribe()
	require.NoError(t, err)

	s.Cancel()

	push(t, l, "first", V(1), V(0), noError)

	v, ok := <-s.C()
	require.False(t, ok)
	require.Equal(t, V(0), v)
}

// TestScan tests Scan
func TestScan(t *testing.T) {
	l := eventlog.New()
	push(t, l, "first", V(1), V(0), noError)
	push(t, l, "second", V(2), V(1), noError)
	push(t, l, "third", V(3), V(2), noError)

	var events []eventlog.Event

	buf := make([]eventlog.Event, 1)
	v, err := eventlog.Scan(l, V(0), buf, func(e eventlog.Event) bool {
		events = append(events, e)
		return true
	})
	require.NoError(t, err)
	require.Equal(t, V(3), v)

	checkExpected(t, []string{"first", "second", "third"}, events, nil)
}

// TestScanInterrupt assumes Scan to return the latest scanned version
// when interrupted
func TestScanInterrupt(t *testing.T) {
	l := eventlog.New()
	push(t, l, "first", V(1), V(0), noError)
	push(t, l, "second", V(2), V(1), noError)
	push(t, l, "third", V(3), V(2), noError)
	push(t, l, "fourth", V(4), V(3), noError)

	for _, t1 := range []struct {
		stopAt        string
		expectVersion V
		expectEvents  []string
	}{
		{"first", V(1), []string{}},
		{"second", V(2), []string{"first"}},
		{"third", V(3), []string{"first", "second"}},
		{"fourth", V(4), []string{"first", "second", "third"}},
	} {
		t.Run(t1.stopAt, func(t *testing.T) {
			var events []eventlog.Event

			buf := make([]eventlog.Event, 1)
			v, err := eventlog.Scan(l, V(0), buf, func(e eventlog.Event) bool {
				if e.Payload.(*testEvent).Name == t1.stopAt {
					return false
				}
				events = append(events, e)
				return true
			})
			require.NoError(t, err)
			require.Equal(t, t1.expectVersion, v)
			checkExpected(t, t1.expectEvents, events, nil)
		})
	}
}

// TestScanAllocBuffer tests Scan assuming it to allocate a buffer
// when nil is passed
func TestScanAllocBuffer(t *testing.T) {
	l := eventlog.New()
	push(t, l, "first", V(1), V(0), noError)
	push(t, l, "second", V(2), V(1), noError)
	push(t, l, "third", V(3), V(2), noError)

	var events []eventlog.Event
	v, err := eventlog.Scan(l, V(0), nil, func(e eventlog.Event) bool {
		events = append(events, e)
		return true
	})
	require.NoError(t, err)
	require.Equal(t, V(3), v)

	checkExpected(t, []string{"first", "second", "third"}, events, nil)
}

type testEvent struct {
	Name string
}

// Copy implements the eventlog.Payload interface
func (e *testEvent) Copy() eventlog.Payload {
	cp := *e
	return &cp
}

var (
	noError  error = nil
	noAssume       = noAssumeT{}
)

type (
	noAssumeT struct{}
	V         = eventlog.Version
)

// push helps pushing events onto the given log
func push(
	t *testing.T,
	l *eventlog.EventLog,
	eventName string,
	expectedNewVersion V,
	assumedVersion interface{},
	expectedError error,
) {
	var (
		initialVersion = l.Version()
		newVersion     V
		pushedEvent    eventlog.Event
		err            error
		newTestEvent   = &testEvent{eventName}
	)
	switch v := assumedVersion.(type) {
	case V:
		newVersion, pushedEvent, err = l.CheckPush(v, newTestEvent)
	case noAssumeT:
		newVersion, pushedEvent, err = l.Push(newTestEvent)
	default:
		t.Fatalf(
			"unexpected input type (%s) for assumedVersion",
			reflect.TypeOf(assumedVersion),
		)
	}
	if expectedError != nil {
		require.Error(t, err)
		require.True(t, errors.Is(err, expectedError))
		require.Equal(t, l.Version(), initialVersion)
		require.Zero(t, pushedEvent)
	} else {
		require.NoError(t, err)
		require.WithinDuration(
			t,
			time.Now().UTC(),
			pushedEvent.Time,
			MaxTimeDelta,
		)
		require.Equal(t, expectedNewVersion, newVersion)
		require.Greater(t, uint64(l.Version()), uint64(initialVersion))
		require.NotZero(t, pushedEvent.ID)
		require.Equal(t, newTestEvent, pushedEvent.Payload)
	}
}

// MaxTimeDelta defines the maximum difference between the time
// of the created event and the time of checking
const MaxTimeDelta = time.Second

// checkExpected helps comparing expected events against the actual ones
func checkExpected(
	t *testing.T,
	expected []string,
	actual []eventlog.Event,
	idRegister map[eventlog.EventID]struct{},
) (newIDRegister map[eventlog.EventID]struct{}) {
	require.Equal(t, len(expected), len(actual))

	newIDRegister = idRegister
	if newIDRegister == nil {
		newIDRegister = make(map[eventlog.EventID]struct{}, len(actual))
	}

	for i, name := range expected {
		e := actual[i]
		require.NotNil(t, e.Payload)
		require.IsType(t, &testEvent{}, e.Payload)
		require.Equal(t, name, e.Payload.(*testEvent).Name)

		require.NotZero(t, e.ID)
		require.NotZero(t, e.Time)

		require.NotContains(t, newIDRegister, e.ID)
		newIDRegister[e.ID] = struct{}{}

		if i > 0 {
			require.GreaterOrEqual(
				t,
				e.Time.Unix(),
				actual[i-1].Time.Unix(),
			)
		}
	}

	return newIDRegister
}
