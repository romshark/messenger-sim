package graph

import (
	"simulator/service/auth"
	"simulator/service/messaging"
	"simulator/service/users"
)

// Resolver represents the graph resolver
type Resolver struct {
	UsersService     users.Client
	MessagingService messaging.Client
	AuthService      auth.Client
}
