package client

import (
	"context"

	"github.com/romshark/messenger-sim/messenger/event"
)

// LeaveConversation makes the logged in user leave the given conversation
func (c *Client) LeaveConversation(
	ctx context.Context,
	conversationID event.ConversationID,
) error {
	return c.Query(
		ctx,
		nil,
		`mutation (
			$conversationID: String!,
		) {
			leaveConversation(
				conversationID: $conversationID,
			)
		}`,
		Args{"conversationID": conversationID},
	)
}
