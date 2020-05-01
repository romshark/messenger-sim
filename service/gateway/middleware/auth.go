package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/messenger/sessid"
	"github.com/romshark/messenger-sim/service/auth"
)

// Auth represents the authentication middleware
type Auth struct {
	logger     Logger
	next       http.Handler
	authClient auth.Client
}

// Logger represents any logger implementation
type Logger interface {
	Printf(format string, a ...interface{}) (int, error)
}

// ConsoleLog is a fallback noop logger
type ConsoleLog struct{}

// Printf implements interface Logger
func (ConsoleLog) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Printf(format, a...)
}

// NewAuth creates a new authentication middleware instance
func NewAuth(
	next http.Handler,
	logger Logger,
	authClient auth.Client,
) (*Auth, error) {
	if next == nil {
		return nil, fmt.Errorf("missing next")
	}

	if authClient == nil {
		return nil, fmt.Errorf("missing authentication service client")
	}

	if logger == nil {
		logger = ConsoleLog{}
	}

	return &Auth{
		logger:     logger,
		next:       next,
		authClient: authClient,
	}, nil
}

func (a *Auth) ServeHTTP(
	resp http.ResponseWriter,
	req *http.Request,
) {
	r := &Request{
		IP:             getIP(req),
		UserAgent:      req.Header.Get("User-Agent"),
		ResponseWriter: resp,
	}

	if c, err := req.Cookie(CookieSessionID); err == nil {
		session, err := a.authClient.FindSessionByID(
			req.Context(),
			sessid.SessionID(c.Value),
		)
		if err != nil {
			a.logger.Printf("searching for session by ID: %s", err)
		} else {
			r.Session = session
		}
	}

	r.Request = req.WithContext(
		context.WithValue(req.Context(), CtxRequest, r),
	)

	a.next.ServeHTTP(resp, r.Request)
}

// getIP extracts the IP address by reading off the forwarded-for
// header (for proxies) and falls back to the remote address
func getIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-FORWARDED-FOR"); forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// Either succeeds if either of the rules can applicable
func Either(
	ctx context.Context,
	rules ...Authorizer,
) error {
	var s *auth.Session
	if req, ok := ctx.Value(CtxRequest).(*Request); ok {
		s = req.Session
	}

	for _, r := range rules {
		if r.Authorize(s) {
			return nil
		}
	}
	return ErrUnauthorized
}

// Authorizer represents an abstract authorization rule
type Authorizer interface {
	Authorize(session *auth.Session) bool
}

// Owner represents an authorization rule assuming the requesting
// user to be the owner of the requested resource
type Owner struct{ ID event.UserID }

// Authorize implements the Authorizer interface
func (o Owner) Authorize(session *auth.Session) bool {
	return session != nil && session.User == o.ID
}

// Authenticated represents an authorization rule
// assuming the requesting user to be authenticated
type Authenticated struct{}

// Authorize implements the Authorizer interface
func (o Authenticated) Authorize(session *auth.Session) bool {
	return session != nil
}

// ErrUnauthorized indicates insufficient permissions
var ErrUnauthorized = errors.New("ErrUnauthorized")

// CtxKey represents the context key type
type CtxKey int

// CtxRequest defines the context key for the request object
const CtxRequest CtxKey = 1

// Request represents the aggregated request information
type Request struct {
	IP             string
	UserAgent      string
	Session        *auth.Session
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// CookieSessionID defines the name of the cookie carrying the session ID
const CookieSessionID = "SID"
