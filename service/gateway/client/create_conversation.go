package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/service/gateway/client/model"
)

// CreateConversation creates a new conversation
func (c *Client) CreateConversation(
	ctx context.Context,
	title string,
	participants []event.UserID,
	avatarURL *url.URL,
	resultQuery string,
) (*model.Conversation, error) {
	var resp struct {
		X *model.Conversation `json:"createConversation"`
	}

	if err := c.Query(
		ctx,
		&resp,
		fmt.Sprintf(`mutation (
			$title: String!,
			$participants: [String!]!,
			$avatarURL: String,
		) {
			createConversation(
				title: $title,
				participants: $participants,
				avatarURL: $avatarURL,
			) { %s }
		}`, resultQuery),
		Args{
			"title":        title,
			"participants": participants,
			"avatarURL":    avatarURL,
		},
	); err != nil {
		return nil, err
	}
	return resp.X, nil
}
