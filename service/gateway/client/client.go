package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/romshark/messenger-sim/service/gateway/client/model"
)

// Client represents an HTTP GraphQL client
type Client struct {
	httpClient *http.Client
	sessionID  string
	userID     string
	endpoint   string
}

// NewClient creates a new client instance
func NewClient(
	httpClient *http.Client,
	endpoint string,
) *Client {
	return &Client{
		httpClient: httpClient,
		endpoint:   endpoint,
	}
}

// SessionID returns the ID of the currently active session
func (c *Client) SessionID() string { return c.sessionID }

// UserID returns the user ID of the currently active session
func (c *Client) UserID() string { return c.userID }

// Query performs a read-only query against the graph
func (c *Client) Query(
	ctx context.Context,
	response interface{},
	query string,
	arguments Args,
) error {
	query = strings.ReplaceAll(query, "\n", " ")
	query = strings.ReplaceAll(query, "\t", " ")

	reqJSONBytes, err := json.Marshal(GQLRequest{
		Query:     query,
		Variables: arguments,
	})
	if err != nil {
		return fmt.Errorf("encoding request JSON: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		c.endpoint,
		bytes.NewBuffer(reqJSONBytes),
	)
	if err != nil {
		return fmt.Errorf("initializing HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("performing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"unexpected response status (%d): %s",
			resp.StatusCode,
			resp.Status,
		)
	}

	respData := GQLResponse{Data: response}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return fmt.Errorf("decoding JSON response: %w", err)
	}

	switch len(respData.Errors) {
	case 0:
	case 1:
		return respData.Errors[0]
	default:
		return fmt.Errorf("errors: %v", respData.Errors)
	}

	return nil
}

// Auth authenticates the client
func (c *Client) Auth(
	ctx context.Context,
	username, password string,
) error {
	var r struct {
		CreateSession *model.Session `json:"createSession"`
	}
	if err := c.Query(
		ctx,
		&r,
		`mutation ($u: String!, $p: String!) {
			createSession(username: $u, password: $p) {
				id
				user { id }
			}
		}`,
		Args{"u": username, "p": password},
	); err != nil {
		return err
	}
	c.sessionID = r.CreateSession.ID
	c.userID = r.CreateSession.User.ID
	return nil
}

// Args is used for query arguments
type Args map[string]interface{}

// GQLRequest represents a GraphQL JSON request
type GQLRequest struct {
	Query         string `json:"query,omitempty"`
	OperationName string `json:"operationName,omitempty"`
	Variables     Args   `json:"variables,omitempty"`
}

// GQLResponse represents a GraphQL JSON response
type GQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []GQLError  `json:"errors"`
}

// GQLError represents a GraphQL error
type GQLError struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

func (e GQLError) Error() string {
	return e.Message
}
