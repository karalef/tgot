package tgot

import (
	"errors"
	"sync/atomic"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// NewWithToken creates new bit with specified token and default http client.
func NewWithToken(token string) (*Bot, error) {
	a, err := api.New(token, "", "", nil)
	if err != nil {
		return nil, err
	}
	return New(a)
}

// New creates new bot.
func New(api *api.API) (*Bot, error) {
	if api == nil {
		return nil, errors.New("nil api")
	}
	b := &Bot{
		api: api,
	}
	_, err := b.GetMe()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Bot type.
type Bot struct {
	err atomic.Pointer[error]
	api *api.API
	me  tg.User

	// OnError is called when an error occurs while method execution with context.
	// The returned error will be returned to the method caller.
	// If this function is nil, OnErrorDefault is used.
	OnError func(c Context, result any, err error) error

	Handler
}

// API returns api object.
func (b *Bot) API() *api.API {
	return b.api
}

// Me returns current bot as tg.User.
func (b *Bot) Me() tg.User {
	return b.me
}

// GetMe returns basic information about the bot in form of a User object.
func (b *Bot) GetMe() (*tg.User, error) {
	me, err := api.Request[*tg.User](b.api, "getMe", nil)
	if err != nil {
		return nil, err
	}
	b.me = *me
	return me, nil
}

// LogOut method.
//
// Use this method to log out from the cloud Bot API server before launching the bot locally.
func (b *Bot) LogOut() error {
	return b.api.Request("logOut", nil)
}

// Close method.
//
// Use this method to close the bot instance before moving it from one local server to another.
func (b *Bot) Close() error {
	return b.api.Request("close", nil)
}

// SetError changes the state of the bot to an error.
// It does nothing if the bot is already in an error state.
// The error state prevents any updates handling.
func (b *Bot) SetError(err error) {
	b.err.CompareAndSwap(nil, &err)
}

// Err returns bot error.
func (b *Bot) Err() error {
	if e := b.err.Load(); e != nil {
		return *e
	}
	return nil
}

func (b *Bot) onError(c Context, result any, err error) error {
	fn := b.OnError
	if fn == nil {
		fn = OnErrorDefault
	}
	return fn(c, result, err)
}

// Handle routes the update to the matching handler.
func (b *Bot) Handle(upd *tg.Update) error {
	if err := b.Err(); err != nil {
		return err
	}

	switch {
	default:
		return ErrNotSpecified

	case upd.Message != nil && b.OnMessage != nil:
		ctx := b.makeChatContext(upd.Message.Chat, "Message")
		b.OnMessage(ctx, upd.Message)

	case upd.EditedMessage != nil && b.OnEditedMessage != nil:
		ctx := b.makeChatContext(upd.EditedMessage.Chat, "EditedMessage")
		b.OnEditedMessage(ctx, upd.EditedMessage)

	case upd.ChannelPost != nil && b.OnChannelPost != nil:
		ctx := b.makeChatContext(upd.ChannelPost.Chat, "Post")
		b.OnChannelPost(ctx, upd.ChannelPost)

	case upd.EditedChannelPost != nil && b.OnEditedChannelPost != nil:
		ctx := b.makeChatContext(upd.EditedChannelPost.Chat, "EditedPost")
		b.OnEditedChannelPost(ctx, upd.EditedChannelPost)

	case upd.CallbackQuery != nil && b.OnCallbackQuery != nil:
		ctx := makeQueryContext[CallbackAnswer](b.MakeContext("Callback"), upd.CallbackQuery.ID)
		b.OnCallbackQuery(ctx, upd.CallbackQuery)

	case upd.InlineQuery != nil && b.OnInlineQuery != nil:
		ctx := makeQueryContext[InlineAnswer](b.MakeContext("Inline"), upd.InlineQuery.ID)
		b.OnInlineQuery(ctx, upd.InlineQuery)

	case upd.InlineChosen != nil && b.OnInlineChosen != nil:
		ctx := b.MakeContext("InlineChosen").OpenMessage(InlineMsgID(upd.InlineChosen))
		b.OnInlineChosen(ctx, upd.InlineChosen)

	case upd.ShippingQuery != nil && b.OnShippingQuery != nil:
		ctx := makeQueryContext[ShippingAnswer](b.MakeContext("ShippingQuery"), upd.ShippingQuery.ID)
		b.OnShippingQuery(ctx, upd.ShippingQuery)

	case upd.PreCheckoutQuery != nil && b.OnPreCheckoutQuery != nil:
		ctx := makeQueryContext[PreCheckoutAnswer](b.MakeContext("PreCheckoutQuery"), upd.PreCheckoutQuery.ID)
		b.OnPreCheckoutQuery(ctx, upd.PreCheckoutQuery)

	case upd.Poll != nil && b.OnPoll != nil:
		b.OnPoll(b.MakeContext("Poll"), upd.Poll)

	case upd.PollAnswer != nil && b.OnPollAnswer != nil:
		b.OnPollAnswer(b.MakeContext("PollAnswer"), upd.PollAnswer)

	case upd.MyChatMember != nil && b.OnMyChatMember != nil:
		ctx := b.makeChatContext(upd.MyChatMember.Chat, "MyChatMember")
		b.OnMyChatMember(ctx, upd.MyChatMember)

	case upd.ChatMember != nil && b.OnChatMember != nil:
		ctx := b.makeChatContext(upd.ChatMember.Chat, "ChatMember")
		b.OnChatMember(ctx, upd.ChatMember)

	case upd.ChatJoinRequest != nil && b.OnChatJoinRequest != nil:
		ctx := b.makeChatContext(upd.ChatJoinRequest.Chat, "JoinRequest")
		b.OnChatJoinRequest(ctx, upd.ChatJoinRequest)
	}

	return b.Err()
}

// ErrNotSpecified is returned by Handle if the update contains an object not specified in Handler.
var ErrNotSpecified = errors.New("handler is not specified")

// Handler contains all available updates handler functions.
type Handler struct {
	OnMessage           func(ChatContext, *tg.Message)
	OnEditedMessage     func(ChatContext, *tg.Message)
	OnChannelPost       func(ChatContext, *tg.Message)
	OnEditedChannelPost func(ChatContext, *tg.Message)
	OnCallbackQuery     func(CallbackContext, *tg.CallbackQuery)
	OnInlineQuery       func(InlineContext, *tg.InlineQuery)
	OnInlineChosen      func(MessageContext, *tg.InlineChosen)
	OnShippingQuery     func(ShippingContext, *tg.ShippingQuery)
	OnPreCheckoutQuery  func(PreCheckoutContext, *tg.PreCheckoutQuery)
	OnPoll              func(Context, *tg.Poll)
	OnPollAnswer        func(Context, *tg.PollAnswer)
	OnMyChatMember      func(ChatContext, *tg.ChatMemberUpdated)
	OnChatMember        func(ChatContext, *tg.ChatMemberUpdated)
	OnChatJoinRequest   func(ChatContext, *tg.ChatJoinRequest)
}

// Allowed returns list of allowed updates.
// If there is any change in the handler, this function must be called to get a new list,
// otherwise it may cause a panic.
func (h *Handler) Allowed() []string {
	list := make([]string, 0, 14)
	add := func(a bool, s string) {
		if a {
			list = append(list, s)
		}
	}
	add(h.OnMessage != nil, "message")
	add(h.OnEditedMessage != nil, "edited_message")
	add(h.OnChannelPost != nil, "channel_post")
	add(h.OnEditedChannelPost != nil, "edited_channel_post")
	add(h.OnCallbackQuery != nil, "callback_query")
	add(h.OnInlineQuery != nil, "inline_query")
	add(h.OnInlineChosen != nil, "chosen_inline_result")
	add(h.OnShippingQuery != nil, "shipping_query")
	add(h.OnPreCheckoutQuery != nil, "pre_checkout_query")
	add(h.OnPoll != nil, "poll")
	add(h.OnPollAnswer != nil, "poll_answer")
	add(h.OnMyChatMember != nil, "my_chat_member")
	add(h.OnChatMember != nil, "chat_member")
	add(h.OnChatJoinRequest != nil, "chat_join_request")
	return list
}
