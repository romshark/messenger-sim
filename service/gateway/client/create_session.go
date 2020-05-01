package client

import (
	"context"
	"fmt"

	"github.com/romshark/messenger-sim/service/gateway/client/model"
)

// CreateSession creates a new session
// automatically authenticating the client with it
func (c *Client) CreateSession(
	ctx context.Context,
	username string,
	password string,
	resultQuery string,
) (*model.Session, error) {
	var resp struct {
		X *model.Session `json:"createSession"`
	}
	if err := c.Query(
		ctx,
		&resp,
		fmt.Sprintf(`mutation (
			$username: String!,
			$password: String!,
		) {
			createSession(
				username:$username,
				password:$password,
			) {
				id
				user { id }
				%s
			}
		}`, resultQuery),
		Args{
			"username": username,
			"password": password,
		},
	); err != nil {
		return nil, err
	}
	c.sessionID = resp.X.ID
	c.userID = resp.X.User.ID
	return resp.X, nil
}
