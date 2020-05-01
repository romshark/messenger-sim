package client

import (
	"context"
	"fmt"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/service/gateway/client/model"
)

// EditMessage edits an existing message
func (c *Client) EditMessage(
	ctx context.Context,
	messageID event.MessageID,
	body string,
	resultQuery string,
) (*model.Message, error) {
	var resp struct {
		X *model.Message `json:"editMessage"`
	}
	if err := c.Query(
		ctx,
		&resp,
		fmt.Sprintf(`mutation (
			$messageID: String!,
			$body: String!,
		) {
			editMessage(
				messageID: $messageID,
				body: $body,
			) { %s }
		}`, resultQuery),
		Args{
			"messageID": messageID,
			"body":      body,
		},
	); err != nil {
		return nil, err
	}
	return resp.X, nil
}
