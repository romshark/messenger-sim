package simulator

import (
	"simulator/messenger/event"
	"simulator/messenger/eventlog"
	"simulator/service/users"
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
			}
			return true
		},
	)
	s.projectionVersion = latestVersion
	return latestVersion, err
}

func (s *Simulator) applyUserCreated(e eventlog.Event, p *event.UserCreated) {
	newUser := &users.User{
		ID:           p.ID,
		Username:     p.Username,
		DisplayName:  p.DisplayName,
		CreationTime: e.Time,
	}
	s.usersByUsername[p.Username] = newUser
	s.usersByID[p.ID] = newUser
}
