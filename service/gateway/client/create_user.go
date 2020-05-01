package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/romshark/messenger-sim/messenger/username"
	"github.com/romshark/messenger-sim/service/gateway/client/model"
)

// CreateUser creates a new user account
func (c *Client) CreateUser(
	ctx context.Context,
	username username.Username,
	displayName string,
	avatarURL *url.URL,
	password string,
	resultQuery string,
) (*model.User, error) {
	if err := username.Validate(); err != nil {
		return nil, fmt.Errorf("invalid username: %w", err)
	}

	var resp struct {
		X *model.User `json:"createUser"`
	}
	if err := c.Query(
		ctx,
		&resp,
		fmt.Sprintf(`mutation (
			$username: String!,
			$displayName: String!,
			$avatarURL: String,
			$password: String!,
		) {
			createUser(
				username:$username,
				displayName: $displayName,
				avatarURL:$avatarURL,
				password:$password,
			) { %s }
		}`, resultQuery),
		Args{
			"username":    username,
			"displayName": displayName,
			"password":    password,
		},
	); err != nil {
		return nil, err
	}
	return resp.X, nil
}
