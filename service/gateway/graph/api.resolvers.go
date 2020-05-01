package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"net/url"

	"github.com/romshark/messenger-sim/messenger/event"
	libid "github.com/romshark/messenger-sim/messenger/id"
	"github.com/romshark/messenger-sim/messenger/sessid"
	usrname "github.com/romshark/messenger-sim/messenger/username"
	"github.com/romshark/messenger-sim/service/gateway/graph/generated"
	"github.com/romshark/messenger-sim/service/gateway/graph/model"
)

func (r *conversationResolver) Participants(ctx context.Context, obj *model.Conversation) ([]*model.User, error) {
	uids, err := r.MessagingService.ListParticipants(
		ctx,
		obj.ConversationID,
	)
	if err != nil {
		return nil, err
	}

	users, err := r.UsersService.GetUsers(ctx, uids)
	if err != nil {
		return nil, err
	}

	l := make([]*model.User, 0, len(users))
	for _, u := range users {
		var avatarURL *string
		if u.AvatarURL != nil {
			v := u.AvatarURL.String()
			avatarURL = &v
		}

		l = append(l, &model.User{
			UserID:       u.ID,
			ID:           u.ID.String(),
			Username:     string(u.Username),
			DisplayName:  u.DisplayName,
			CreationTime: u.CreationTime,
			AvatarURL:    avatarURL,
		})
	}
	return l, nil
}

func (r *conversationResolver) Messages(ctx context.Context, obj *model.Conversation, afterID *string, limit int) ([]*model.Message, error) {
	var after *event.MessageID
	if afterID != nil {
		id, err := libid.FromString(*afterID)
		if err != nil {
			return nil, fmt.Errorf("parsing after ID: %w", err)
		}
		mid := event.MessageID(id)
		after = &mid
	}

	messages, err := r.MessagingService.GetMessages(
		ctx,
		obj.ConversationID,
		after,
		limit,
	)
	if err != nil {
		return nil, err
	}

	resolvers := make([]*model.Message, len(messages))
	for i, m := range messages {
		resolvers[i] = &model.Message{
			MessageID:      m.ID,
			ID:             m.ID.String(),
			Body:           m.Body,
			SendingTime:    m.SendingTime,
			SenderID:       m.Sender,
			ConversationID: obj.ConversationID,
		}
	}
	return resolvers, nil
}

func (r *messageResolver) Sender(ctx context.Context, obj *model.Message) (*model.User, error) {
	users, err := r.UsersService.GetUsers(ctx, []event.UserID{obj.SenderID})
	if err != nil {
		return nil, err
	}
	if len(users) < 0 {
		return nil, nil
	}
	u := users[0]

	var avatarURL *string
	if u.AvatarURL != nil {
		v := u.AvatarURL.String()
		avatarURL = &v
	}

	return &model.User{
		UserID:       u.ID,
		ID:           u.ID.String(),
		Username:     string(u.Username),
		DisplayName:  u.DisplayName,
		CreationTime: u.CreationTime,
		AvatarURL:    avatarURL,
	}, nil
}

func (r *messageResolver) Conversation(ctx context.Context, obj *model.Message) (*model.Conversation, error) {
	c, err := r.MessagingService.FindConversation(
		ctx,
		obj.ConversationID,
	)
	if err != nil {
		return nil, err
	}

	var avatarURL *string
	if c.AvatarURL != nil {
		v := c.AvatarURL.String()
		avatarURL = &v
	}

	return &model.Conversation{
		ConversationID: c.ID,
		ID:             c.ID.String(),
		Title:          c.Title,
		AvatarURL:      avatarURL,
		CreationTime:   c.CreationTime,
	}, nil
}

func (r *messageResolver) Edits(ctx context.Context, obj *model.Message) ([]*model.MessageEdit, error) {
	e, err := r.MessagingService.GetMessageEdits(ctx, obj.MessageID)
	if err != nil {
		return nil, err
	}

	l := make([]*model.MessageEdit, len(e))
	for i, e := range e {
		l[i] = &model.MessageEdit{
			Time:         e.Time,
			PreviousBody: e.PreviousBody,
		}
	}
	return l, nil
}

func (r *messageEditResolver) Editor(ctx context.Context, obj *model.MessageEdit) (*model.User, error) {
	users, err := r.UsersService.GetUsers(ctx, []event.UserID{obj.EditorID})
	if err != nil {
		return nil, err
	}
	if len(users) < 0 {
		return nil, nil
	}
	u := users[0]

	var avatarURL *string
	if u.AvatarURL != nil {
		v := u.AvatarURL.String()
		avatarURL = &v
	}

	return &model.User{
		UserID:       obj.EditorID,
		ID:           obj.EditorID.String(),
		Username:     string(u.Username),
		DisplayName:  u.DisplayName,
		CreationTime: u.CreationTime,
		AvatarURL:    avatarURL,
	}, nil
}

func (r *mutationResolver) CreateSession(ctx context.Context, username string, password string) (*model.Session, error) {
	// TODO: get IP & UserAgent from context
	const nullIP = "0.0.0.0"
	sess, err := r.AuthService.CreateSession(
		ctx,
		usrname.Username(username),
		password,
		nullIP,
		"",
	)
	if err != nil {
		return nil, err
	}
	return &model.Session{
		ID:           string(sess.ID),
		IP:           nullIP,
		UserAgent:    "",
		CreationTime: sess.CreationTime,
	}, nil
}

func (r *mutationResolver) DestroySession(ctx context.Context, id string) (bool, error) {
	if err := r.AuthService.DestroySession(
		ctx,
		sessid.SessionID(id),
	); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, username string, displayName string, avatarURL *string, password string) (*model.User, error) {
	var avatar *url.URL
	if avatarURL != nil {
		var err error
		avatar, err = url.Parse(*avatarURL)
		if err != nil {
			return nil, fmt.Errorf("parsing avatar URL: %w", err)
		}
	}

	newUser, err := r.UsersService.CreateNewUser(
		ctx,
		usrname.Username(username),
		displayName,
		avatar,
		password,
	)
	if err != nil {
		return nil, err
	}
	return &model.User{
		UserID:       newUser.ID,
		ID:           newUser.ID.String(),
		Username:     username,
		DisplayName:  displayName,
		CreationTime: newUser.CreationTime,
		AvatarURL:    avatarURL,
	}, nil
}

func (r *mutationResolver) SendMessage(ctx context.Context, body string, conversationID string) (*model.Message, error) {
	convID, err := libid.FromString(conversationID)
	if err != nil {
		return nil, fmt.Errorf("parsing conversation ID: %w", err)
	}

	newMsg, err := r.MessagingService.SendMessage(
		ctx,
		body,
		convID,
		event.UserID{}, // SenderID
	)
	if err != nil {
		return nil, err
	}

	return &model.Message{
		MessageID:      newMsg.ID,
		SenderID:       newMsg.Sender,
		ConversationID: newMsg.Conversation,
		ID:             newMsg.ID.String(),
		Body:           newMsg.Body,
		SendingTime:    newMsg.SendingTime,
	}, nil
}

func (r *mutationResolver) EditMessage(ctx context.Context, messageID string, body string) (*model.Message, error) {
	mid, err := libid.FromString(messageID)
	if err != nil {
		return nil, fmt.Errorf("parsing message ID: %w", err)
	}

	// TODO: get editor ID from session
	var editorID event.UserID

	editedMsg, err := r.MessagingService.EditMessage(
		ctx,
		event.MessageID(mid),
		editorID,
		body,
	)
	if err != nil {
		return nil, err
	}
	return &model.Message{
		MessageID:      event.MessageID(mid),
		SenderID:       editedMsg.Sender,
		ConversationID: editedMsg.Conversation,
		ID:             messageID,
		Body:           editedMsg.Body,
		SendingTime:    editedMsg.SendingTime,
	}, nil
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, messageID string, reason *string) (bool, error) {
	mid, err := libid.FromString(messageID)
	if err != nil {
		return false, fmt.Errorf("parsing message ID: %w", err)
	}

	// TODO: get deletor ID from session
	var deletorID event.UserID

	err = r.MessagingService.DeleteMessage(
		ctx,
		event.MessageID(mid),
		deletorID,
		reason,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) CreateConversation(ctx context.Context, title string, participants []string, avatarURL *string) (*model.Conversation, error) {
	participantIDs := make([]event.UserID, len(participants))
	for i, p := range participants {
		id, err := libid.FromString(p)
		if err != nil {
			return nil, fmt.Errorf(
				"parsing participant ID (%d): %w",
				i, err,
			)
		}
		participantIDs[i] = event.UserID(id)
	}

	var avatar *url.URL
	if avatarURL != nil {
		var err error
		if avatar, err = url.Parse(*avatarURL); err != nil {
			return nil, fmt.Errorf("parsing avatar URL: %w", err)
		}
	}

	// TODO: get creator ID from session
	var creatorID event.UserID

	newConv, err := r.MessagingService.CreateConversation(
		ctx,
		title,
		creatorID,
		participantIDs,
		avatar,
	)
	if err != nil {
		return nil, err
	}

	return &model.Conversation{
		ConversationID: newConv.ID,
		ID:             newConv.ID.String(),
		Title:          title,
		AvatarURL:      avatarURL,
		CreationTime:   newConv.CreationTime,
	}, nil
}

func (r *mutationResolver) LeaveConversation(ctx context.Context, conversationID string) (bool, error) {
	// TODO: get user ID from session
	var userID event.UserID

	cid, err := libid.FromString(conversationID)
	if err != nil {
		return false, fmt.Errorf("parsing conversation ID: %w", err)
	}

	err = r.MessagingService.LeaveConversation(
		ctx,
		userID,
		event.ConversationID(cid),
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) EditConversation(ctx context.Context, conversationID string, title *string, avatarURL *string) (*model.Conversation, error) {
	cid, err := libid.FromString(conversationID)
	if err != nil {
		return nil, fmt.Errorf("parsing conversation ID: %w", err)
	}

	// TODO: get editor ID from session
	var editorID event.UserID

	var avatar interface{}
	if avatarURL != nil {
		if *avatarURL == "" {
			// reset
			var nilPtr *event.UserID
			avatar = nilPtr
		} else {
			var err error
			if avatar, err = url.Parse(*avatarURL); err != nil {
				return nil, fmt.Errorf("parsing avatar URL: %w", err)
			}
		}
	}

	updatedConv, err := r.MessagingService.EditConversation(
		ctx,
		event.ConversationID(cid),
		editorID,
		title,
		avatar,
	)
	if err != nil {
		return nil, err
	}

	var newAvatar *string
	if updatedConv.AvatarURL != nil {
		v := updatedConv.AvatarURL.String()
		newAvatar = &v
	}

	return &model.Conversation{
		ConversationID: updatedConv.ID,
		ID:             conversationID,
		Title:          updatedConv.Title,
		AvatarURL:      avatarURL,
		CreationTime:   updatedConv.CreationTime,
	}, nil
}

func (r *mutationResolver) RemoveUserFromConversation(ctx context.Context, conversationID string, userID string) (bool, error) {
	cid, err := libid.FromString(conversationID)
	if err != nil {
		return false, fmt.Errorf("parsing conversation ID: %w", err)
	}

	uid, err := libid.FromString(userID)
	if err != nil {
		return false, fmt.Errorf("parsing user ID: %w", err)
	}

	// TODO: get remover ID from session
	var removerID event.UserID

	err = r.MessagingService.RemoveUserFromConversation(
		ctx,
		event.ConversationID(cid),
		event.UserID(uid),
		removerID,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	i, err := libid.FromString(id)
	if err != nil {
		return nil, fmt.Errorf("parsing user ID: %w", err)
	}
	uid := event.UserID(i)

	users, err := r.UsersService.GetUsers(ctx, []event.UserID{uid})
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, nil
	}
	u := users[0]

	var avatarURL *string
	if u.AvatarURL != nil {
		v := u.AvatarURL.String()
		avatarURL = &v
	}

	return &model.User{
		UserID:       uid,
		ID:           id,
		Username:     string(u.Username),
		DisplayName:  u.DisplayName,
		AvatarURL:    avatarURL,
		CreationTime: u.CreationTime,
	}, nil
}

func (r *sessionResolver) User(ctx context.Context, obj *model.Session) (*model.User, error) {
	users, err := r.UsersService.GetUsers(ctx, []event.UserID{obj.UserID})
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, nil
	}
	u := users[0]

	var avatarURL *string
	if u.AvatarURL != nil {
		v := u.AvatarURL.String()
		avatarURL = &v
	}

	return &model.User{
		UserID:       obj.UserID,
		ID:           u.ID.String(),
		Username:     string(u.Username),
		DisplayName:  u.DisplayName,
		CreationTime: u.CreationTime,
		AvatarURL:    avatarURL,
	}, nil
}

func (r *userResolver) Sessions(ctx context.Context, obj *model.User) ([]*model.Session, error) {
	sessions, err := r.AuthService.ListSessionsForUser(ctx, obj.UserID)
	if err != nil {
		return nil, err
	}

	l := make([]*model.Session, len(sessions))
	for i, s := range sessions {
		l[i] = &model.Session{
			UserID:       obj.UserID,
			ID:           string(s.ID),
			IP:           s.IP,
			UserAgent:    s.UserAgent,
			CreationTime: s.CreationTime,
		}
	}
	return l, nil
}

func (r *userResolver) Conversations(ctx context.Context, obj *model.User) ([]*model.Conversation, error) {
	convs, err := r.MessagingService.ListConversationsForUser(ctx, obj.UserID)
	if err != nil {
		return nil, err
	}

	l := make([]*model.Conversation, 0, len(convs))
	for _, c := range convs {
		var avatar *string
		if c.AvatarURL != nil {
			v := c.AvatarURL.String()
			avatar = &v
		}

		l = append(l, &model.Conversation{
			ConversationID: c.ID,
			ID:             c.ID.String(),
			Title:          c.Title,
			AvatarURL:      avatar,
			CreationTime:   c.CreationTime,
		})
	}
	return l, nil
}

// Conversation returns generated.ConversationResolver implementation.
func (r *Resolver) Conversation() generated.ConversationResolver { return &conversationResolver{r} }

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// MessageEdit returns generated.MessageEditResolver implementation.
func (r *Resolver) MessageEdit() generated.MessageEditResolver { return &messageEditResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Session returns generated.SessionResolver implementation.
func (r *Resolver) Session() generated.SessionResolver { return &sessionResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type (
	conversationResolver struct{ *Resolver }
	messageResolver      struct{ *Resolver }
	messageEditResolver  struct{ *Resolver }
	mutationResolver     struct{ *Resolver }
	queryResolver        struct{ *Resolver }
	sessionResolver      struct{ *Resolver }
	userResolver         struct{ *Resolver }
)
