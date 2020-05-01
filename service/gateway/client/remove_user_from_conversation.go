package client

import (
	"context"

	"github.com/romshark/messenger-sim/messenger/event"
)

// RemoveUserFromConversation removes a user from the given conversation
func (c *Client) RemoveUserFromConversation(
	ctx context.Context,
	conversationID event.ConversationID,
	userID event.UserID,
) error {
	return c.Query(
		ctx,
		nil,
		`mutation (
			$conversationID: String!,
			$userID: String!,
		) {
			removeUserFromConversation(
				conversationID: $conversationID,
				userID: $userID,
			)
		}`,
		Args{
			"conversationID": conversationID,
			"userID":         userID,
		},
	)
}
