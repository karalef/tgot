package tgot

import (
	"errors"
	"sync/atomic"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/logger"
)

// NewWithToken creates new bit with specified token and default http client.
func NewWithToken(token string, log ...*logger.Logger) (*Bot, error) {
	a, err := api.New(token, "", "", nil)
	if err != nil {
		return nil, err
	}
	return New(a, log...)
}

// New creates new bot.
func New(api *api.API, log ...*logger.Logger) (*Bot, error) {
	if api == nil {
		return nil, errors.New("nil api")
	}
	me, err := api.GetMe()
	if err != nil {
		return nil, err
	}

	b := &Bot{
		api: api,
		me:  *me,
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
	err atomic.Pointer[error]
	api *api.API
	log *logger.Logger
	me  tg.User

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

// Handle routes the update to the matching handler.
// It panics if the update contains an object not specified in bot.Allowed.
func (b *Bot) Handle(upd *tg.Update) error {
	if err := b.Err(); err != nil {
		return err
	}

	switch {
	case upd.Message != nil:
		ctx := b.makeChatContext(upd.Message.Chat, "Message")
		b.OnMessage(ctx, upd.Message)
	case upd.EditedMessage != nil:
		ctx := b.makeChatContext(upd.EditedMessage.Chat, "EditedMessage")
		b.OnEditedMessage(ctx, upd.EditedMessage)
	case upd.ChannelPost != nil:
		ctx := b.makeChatContext(upd.ChannelPost.Chat, "Post")
		b.OnChannelPost(ctx, upd.ChannelPost)
	case upd.EditedChannelPost != nil:
		ctx := b.makeChatContext(upd.EditedChannelPost.Chat, "EditedPost")
		b.OnEditedChannelPost(ctx, upd.EditedChannelPost)
	case upd.CallbackQuery != nil:
		ctx := makeQueryContext[CallbackAnswer](b.MakeContext("Callback"), upd.CallbackQuery.ID)
		b.OnCallbackQuery(ctx, upd.CallbackQuery)
	case upd.InlineQuery != nil:
		ctx := makeQueryContext[InlineAnswer](b.MakeContext("Inline"), upd.InlineQuery.ID)
		b.OnInlineQuery(ctx, upd.InlineQuery)
	case upd.InlineChosen != nil:
		ctx := b.MakeContext("InlineChosen").OpenMessage(InlineSignature(upd.InlineChosen))
		b.OnInlineChosen(ctx, upd.InlineChosen)
	case upd.ShippingQuery != nil:
		ctx := makeQueryContext[ShippingAnswer](b.MakeContext("ShippingQuery"), upd.ShippingQuery.ID)
		b.OnShippingQuery(ctx, upd.ShippingQuery)
	case upd.PreCheckoutQuery != nil:
		ctx := makeQueryContext[PreCheckoutAnswer](b.MakeContext("PreCheckoutQuery"), upd.PreCheckoutQuery.ID)
		b.OnPreCheckoutQuery(ctx, upd.PreCheckoutQuery)
	case upd.Poll != nil:
		b.OnPoll(b.MakeContext("Poll"), upd.Poll)
	case upd.PollAnswer != nil:
		b.OnPollAnswer(b.MakeContext("PollAnswer"), upd.PollAnswer)
	case upd.MyChatMember != nil:
		ctx := b.makeChatContext(upd.MyChatMember.Chat, "MyChatMember")
		b.OnMyChatMember(ctx, upd.MyChatMember)
	case upd.ChatMember != nil:
		ctx := b.makeChatContext(upd.ChatMember.Chat, "ChatMember")
		b.OnChatMember(ctx, upd.ChatMember)
	case upd.ChatJoinRequest != nil:
		ctx := b.makeChatContext(upd.ChatJoinRequest.Chat, "JoinRequest")
		b.OnChatJoinRequest(ctx, upd.ChatJoinRequest)
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
