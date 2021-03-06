package gateway

import (
	"fmt"
	"net/http"

	"github.com/romshark/messenger-sim/service/auth"
	"github.com/romshark/messenger-sim/service/gateway/graph"
	"github.com/romshark/messenger-sim/service/gateway/graph/generated"
	"github.com/romshark/messenger-sim/service/gateway/middleware"
	"github.com/romshark/messenger-sim/service/messaging"
	"github.com/romshark/messenger-sim/service/users"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
)

// Server represents the gateway ingress server
type Server struct {
	authMiddleware *middleware.Auth
	graphServer    *handler.Server
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

	graphServer := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			UsersService:     usersClient,
			MessagingService: messagingClient,
			AuthService:      authClient,
		},
	}))
	graphServer.Use(extension.Introspection{})

	// graphServer.AddTransport(transport.Websocket{
	// 	KeepAlivePingInterval: 10 * time.Second,
	// })
	// graphServer.AddTransport(transport.GET{})
	graphServer.AddTransport(transport.Options{})
	graphServer.AddTransport(transport.POST{})
	graphServer.AddTransport(transport.MultipartForm{})

	// graphServer.SetQueryCache(lru.New(1000))

	graphServer.Use(extension.Introspection{})
	// graphServer.Use(extension.AutomaticPersistedQuery{
	// 	Cache: lru.New(100),
	// })

	authMiddleware, err := middleware.NewAuth(
		graphServer,
		nil,
		authClient,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"initializing authentication middleware: %w",
			err,
		)
	}

	return &Server{
		authMiddleware: authMiddleware,
		graphServer:    graphServer,
	}, nil
}

func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		switch req.URL.Path {
		case "/g":
			s.authMiddleware.ServeHTTP(resp, req)
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
