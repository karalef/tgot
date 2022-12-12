package tgot

import (
	"errors"
	"sync/atomic"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/logger"
)

// NewWithToken creates new bit with specified token and default http client.
func NewWithToken(token string, handler Handler, log ...*logger.Logger) (*Bot, error) {
	a, err := api.New(token, "", "", nil)
	if err != nil {
		return nil, err
	}
	return New(a, handler, log...)
}

// New creates new bot.
func New(api *api.API, handler Handler, log ...*logger.Logger) (*Bot, error) {
	if api == nil {
		return nil, errors.New("nil api")
	}
	me, err := api.GetMe()
	if err != nil {
		return nil, err
	}

	b := &Bot{
		api:     api,
		handler: handler,
		me:      *me,
	}
	if len(log) > 0 {
		b.log = log[0]
	} else {
		b.log = logger.Default("Bot")
	}

	return b, nil
}

// Bot type.
type Bot struct {
	err     atomic.Pointer[error]
	api     *api.API
	log     *logger.Logger
	handler Handler

	me tg.User
}

// API returns api object.
func (b *Bot) API() *api.API {
	return b.api
}

// Me returns current bot as tg.User.
func (b *Bot) Me() tg.User {
	return b.me
}

func (b *Bot) cancel(err error) {
	b.err.CompareAndSwap(nil, &err)
}

// Err returns bot error.
func (b *Bot) Err() error {
	if e := b.err.Load(); e != nil {
		return *e
	}
	return nil
}

// Allowed returns allowed updates list.
func (b *Bot) Allowed() []string {
	return b.handler.allowed()
}

// Handle routes the update to the matching handler.
// It panics if the update contains an object not specified in bot.Allowed.
func (b *Bot) Handle(upd *tg.Update) error {
	if err := b.Err(); err != nil {
		return err
	}

	h := &b.handler
	switch {
	case upd.Message != nil:
		ctx := b.makeChatContext(upd.Message.Chat, "Message")
		h.OnMessage(ctx, upd.Message)
	case upd.EditedMessage != nil:
		ctx := b.makeChatContext(upd.EditedMessage.Chat, "EditedMessage")
		h.OnEditedMessage(ctx, upd.EditedMessage)
	case upd.ChannelPost != nil:
		ctx := b.makeChatContext(upd.ChannelPost.Chat, "Post")
		h.OnChannelPost(ctx, upd.ChannelPost)
	case upd.EditedChannelPost != nil:
		ctx := b.makeChatContext(upd.EditedChannelPost.Chat, "EditedPost")
		h.OnEditedChannelPost(ctx, upd.EditedChannelPost)
	case upd.CallbackQuery != nil:
		ctx := makeQueryContext[CallbackAnswer](b.MakeContext("Callback"), upd.CallbackQuery.ID)
		h.OnCallbackQuery(ctx, upd.CallbackQuery)
	case upd.InlineQuery != nil:
		ctx := makeQueryContext[InlineAnswer](b.MakeContext("Inline"), upd.InlineQuery.ID)
		h.OnInlineQuery(ctx, upd.InlineQuery)
	case upd.InlineChosen != nil:
		ctx := b.MakeContext("InlineChosen").OpenMessage(InlineSignature(upd.InlineChosen))
		h.OnInlineChosen(ctx, upd.InlineChosen)
	case upd.ShippingQuery != nil:
		ctx := makeQueryContext[ShippingAnswer](b.MakeContext("ShippingQuery"), upd.ShippingQuery.ID)
		h.OnShippingQuery(ctx, upd.ShippingQuery)
	case upd.PreCheckoutQuery != nil:
		ctx := makeQueryContext[PreCheckoutAnswer](b.MakeContext("PreCheckoutQuery"), upd.PreCheckoutQuery.ID)
		h.OnPreCheckoutQuery(ctx, upd.PreCheckoutQuery)
	case upd.Poll != nil:
		h.OnPoll(b.MakeContext("Poll"), upd.Poll)
	case upd.PollAnswer != nil:
		h.OnPollAnswer(b.MakeContext("PollAnswer"), upd.PollAnswer)
	case upd.MyChatMember != nil:
		ctx := b.makeChatContext(upd.MyChatMember.Chat, "MyChatMember")
		h.OnMyChatMember(ctx, upd.MyChatMember)
	case upd.ChatMember != nil:
		ctx := b.makeChatContext(upd.ChatMember.Chat, "ChatMember")
		h.OnChatMember(ctx, upd.ChatMember)
	case upd.ChatJoinRequest != nil:
		ctx := b.makeChatContext(upd.ChatJoinRequest.Chat, "JoinRequest")
		h.OnChatJoinRequest(ctx, upd.ChatJoinRequest)
	}

	return b.Err()
}

// Handler conatains all available updates handler functions.
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

func (h *Handler) allowed() []string {
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
	add(h.OnInlineChosen != nil, "inline_choosen")
	add(h.OnShippingQuery != nil, "shipping_query")
	add(h.OnPreCheckoutQuery != nil, "pre_checkout_query")
	add(h.OnPoll != nil, "poll")
	add(h.OnPollAnswer != nil, "poll_answer")
	add(h.OnMyChatMember != nil, "my_chat_member")
	add(h.OnChatMember != nil, "chat_member")
	add(h.OnChatJoinRequest != nil, "chat_join_request")
	return list
}
