package simulator

import (
	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/sessid"
	"github.com/romshark/messenger-sim/service/auth"
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
			case *event.SessionCreated:
				s.applySessionCreated(e, p)
			case *event.SessionDestroyed:
				s.applySessionDestroyed(e, p)
			}
			return true
		},
	)
	s.projectionVersion = latestVersion
	return latestVersion, err
}

func (s *Simulator) applyUserCreated(e eventlog.Event, p *event.UserCreated) {
	s.usersByUsername[p.Username] = &user{
		id:           p.ID,
		passwordHash: p.PasswordHash,
		sessions:     make(map[sessid.SessionID]*session),
	}
}

func (s *Simulator) applySessionCreated(
	e eventlog.Event,
	p *event.SessionCreated,
) {
	user := func() *user {
		for _, u := range s.usersByUsername {
			if u.id == p.User {
				return u
			}
		}
		return nil
	}()
	newSession := &session{
		Session: auth.Session{
			ID:           p.ID,
			User:         user.id,
			IP:           p.IP,
			UserAgent:    p.UserAgent,
			CreationTime: e.Time,
		},
		user: user,
	}
	user.sessions[p.ID] = newSession
	s.sessionsByID[p.ID] = newSession
}

func (s *Simulator) applySessionDestroyed(
	e eventlog.Event,
	p *event.SessionDestroyed,
) {
	sess := s.sessionsByID[p.Session]
	delete(sess.user.sessions, p.Session)
	delete(s.sessionsByID, p.Session)
}
