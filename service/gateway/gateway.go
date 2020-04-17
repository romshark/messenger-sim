package gateway

import (
	"fmt"
	"net/http"

	"github.com/romshark/messenger-sim/service/auth"
	"github.com/romshark/messenger-sim/service/gateway/graph"
	"github.com/romshark/messenger-sim/service/gateway/graph/generated"
	"github.com/romshark/messenger-sim/service/messaging"
	"github.com/romshark/messenger-sim/service/users"

	"github.com/99designs/gqlgen/graphql/handler"
)

// Server represents the gateway ingress server
type Server struct {
	graphServer *handler.Server
}

// NewServer creates a new Gateway ingress server
func NewServer(
	usersClient users.Client,
	authClient auth.Client,
	messagingClient messaging.Client,
) (*Server, error) {
	if usersClient == nil {
		return nil, fmt.Errorf("missing users service client")
	}
	if authClient == nil {
		return nil, fmt.Errorf("missing authentication service client")
	}
	if messagingClient == nil {
		return nil, fmt.Errorf("missing messaging service client")
	}

	return &Server{
		graphServer: handler.NewDefaultServer(
			generated.NewExecutableSchema(
				generated.Config{Resolvers: &graph.Resolver{
					UsersService:     usersClient,
					MessagingService: messagingClient,
					AuthService:      authClient,
				}},
			),
		),
	}, nil
}

func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		switch req.URL.Path {
		case "/g":
			s.graphServer.ServeHTTP(resp, req)
		default:
			http.Error(
				resp,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound,
			)
			return
		}
	default:
		http.Error(
			resp,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}
