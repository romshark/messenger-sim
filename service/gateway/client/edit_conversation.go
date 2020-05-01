package client

import (
	"context"
	"fmt"

	"github.com/romshark/messenger-sim/messenger/event"
	"github.com/romshark/messenger-sim/service/gateway/client/model"
)

// EditConversation edits an existing conversation
func (c *Client) EditConversation(
	ctx context.Context,
	conversationID event.ConversationID,
	title *string,
	avatarURL *string,
	resultQuery string,
) (*model.Conversation, error) {
	var resp struct {
		X *model.Conversation `json:"editConversation"`
	}
	if err := c.Query(
		ctx,
		&resp,
		fmt.Sprintf(`mutation (
			$conversationID: String!,
			$title: String,
			$avatarURL: String,
		) {
			editConversation(
				conversationID: $conversationID,
				title: $title,
				avatarURL: $avatarURL,
			) { %s }
		}`, resultQuery),
		Args{
			"conversationID": conversationID,
			"title":          title,
			"avatarURL":      avatarURL,
		},
	); err != nil {
		return nil, err
	}
	return resp.X, nil
}
