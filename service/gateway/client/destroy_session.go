package client

import "context"

// DestroySession destroys the current session
func (c *Client) DestroySession(ctx context.Context) error {
	if c.sessionID == "" {
		return nil
	}

	if err := c.Query(
		ctx,
		nil,
		`mutation ($id: String!) {
			destroySession(id: $id)
		}`,
		Args{"id": c.sessionID},
	); err != nil {
		return err
	}
	c.sessionID = ""
	c.userID = ""
	return nil
}
