package gateway

import (
	"fmt"
	"net/http"
	"simulator/service/auth"
	"simulator/service/gateway/graph"
	"simulator/service/gateway/graph/generated"
	"simulator/service/messaging"
	"simulator/service/users"

	"github.com/99designs/gqlgen/graphql/handler"
)

// Server represents the gateway ingress server
type Server struct {
	*http.Server
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

	s := &Server{
		graphServer: handler.NewDefaultServer(
			generated.NewExecutableSchema(
				generated.Config{Resolvers: &graph.Resolver{
					UsersService:     usersClient,
					MessagingService: messagingClient,
					AuthService:      authClient,
				}},
			),
		),
	}
	s.Server = &http.Server{
		Handler: s,
	}
	return s, nil
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
