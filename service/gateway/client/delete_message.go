package client

import (
	"context"

	"github.com/romshark/messenger-sim/messenger/event"
)

// DeleteMessage permanently deletes an existing message
func (c *Client) DeleteMessage(
	ctx context.Context,
	messageID event.MessageID,
	reason string,
) error {
	return c.Query(
		ctx,
		nil,
		`mutation (
			$messageID: String!,
			$reason: String!,
		) {
			deleteMessage(
				messageID: $messageID,
				reason: $reason,
			)
		}`,
		Args{
			"messageID": messageID,
			"reason":    reason,
		},
	)
}
