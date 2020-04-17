package graph

import (
	"github.com/romshark/messenger-sim/service/auth"
	"github.com/romshark/messenger-sim/service/messaging"
	"github.com/romshark/messenger-sim/service/users"
)

// Resolver represents the graph resolver
type Resolver struct {
	UsersService     users.Client
	MessagingService messaging.Client
	AuthService      auth.Client
}
