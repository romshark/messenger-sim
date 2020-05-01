package client

import (
	"context"
	"fmt"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/service/gateway/client/model"
)

// SendMessage sends a new message to a conversation
func (c *Client) SendMessage(
	ctx context.Context,
	body string,
	conversationID event.ConversationID,
	resultQuery string,
) (*model.Message, error) {
	var resp struct {
		X *model.Message `json:"sendMessage"`
	}
	if err := c.Query(
		ctx,
		&resp,
		fmt.Sprintf(`mutation (
			$body: String!,
			$conversationID: String!,
		) {
			sendMessage(
				body: $body,
				conversationID: $conversationID,
			) { %s }
		}`, resultQuery),
		Args{
			"body":           body,
			"conversationID": conversationID,
		},
	); err != nil {
		return nil, err
	}
	return resp.X, nil
}
