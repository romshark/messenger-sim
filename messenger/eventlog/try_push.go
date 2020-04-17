package eventlog

import (
	"context"
	"errors"
	"fmt"
)

// TryPush keeps executing fn until either cancelled,
// succeeded (assumed and actual event log versions match)
// or failed due to an error
func TryPush(
	ctx context.Context,
	eventLog *EventLog,
	assumedVersion Version,
	fn func(retries int) (Payload, error),
	sync func() (Version, error),
) (
	pushedEvent Event,
	err error,
) {
	// Reapeat until either cancelled, succeeded or failed
	for i := 0; ; i++ {
		// Check context for cancelation
		if err = ctx.Err(); err != nil {
			return
		}

		var payload Payload
		payload, err = fn(i)
		if err != nil {
			err = fmt.Errorf("executing transaction: %w", err)
			return
		}

		// Try to push a new event onto the event log
		_, pushedEvent, err = eventLog.CheckPush(assumedVersion, payload)
		switch {
		case errors.Is(err, ErrMismatchingVersion):
			// The projection is out of sync, synchronize & repeat
			if assumedVersion, err = sync(); err != nil {
				err = fmt.Errorf("synchronizing: %w", err)
				return
			}
			continue
		case err != nil:
			// Push failed for unexpected reason
			err = fmt.Errorf("pushing event: %w", err)
			return
		}

		// Transaction successfuly committed, synchronize projection
		if _, err = sync(); err != nil {
			err = fmt.Errorf("finalizing: %w", err)
			return
		}
		break
	}
	return
}
