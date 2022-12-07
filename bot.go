package tgot

import (
	"errors"
	"sync/atomic"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/logger"
	"github.com/karalef/tgot/updates"
)

// Config contains bot configuration.
type Config struct {
	Logger   *logger.Logger // logger.Default if empty
	Handler  Handler
	Commands Commands
}

// New creates new bot.
func New(api *api.API, poller updates.Poller, config Config) (*Bot, error) {
	if api == nil {
		return nil, errors.New("nil api")
	}
	if poller == nil {
		return nil, errors.New("nil poller")
	}
	if config.Logger == nil {
		config.Logger = logger.Default("BOT")
	}
	me, err := api.GetMe()
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:     api,
		log:     config.Logger,
		poller:  poller,
		handler: config.Handler,
		cmds:    config.Commands,
		me:      *me,
	}, nil
}

// Bot type.
type Bot struct {
	api *api.API

	log    *logger.Logger
	poller updates.Poller

	err atomic.Pointer[error]

	handler Handler
	cmds    Commands

	me tg.User
}

func (b *Bot) cancel(err error) {
	if b.err.CompareAndSwap(nil, &err) {
		go b.Stop()
	}
}

// Stop stops polling for updates.
func (b *Bot) Stop() {
	b.poller.Close()
}

// Run starts bot.
func (b *Bot) Run() error {
	if b.cmds != nil {
		if err := b.cmds.Init(b.api); err != nil {
			return err
		}
	}

	allowed := b.handler.allowed(b)
	err := b.poller.Run(b.api, b.handle, allowed)
	if e := b.err.Load(); e != nil {
		return *e
	}
	return err
}

func (b *Bot) handle(upd *tg.Update) {
	h := &b.handler

	// since [Handler.allowed] returns a list of the specified handlers
	// it makes no sense to check for nil except for OnMessage and OnChannelPost.
	switch {
	case upd.Message != nil:
		b.onMessage(upd.Message)
	case upd.EditedMessage != nil:
		ctx := b.makeChatContext(upd.EditedMessage.Chat, "EditedMessage")
		h.OnEditedMessage(ctx, upd.EditedMessage)
	case upd.ChannelPost != nil:
		b.onMessage(upd.ChannelPost)
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
}

func (b *Bot) onCommand(msg *tg.Message, cmd string, args []string) bool {
	command := b.cmds.Get(cmd, msg)
	if command == nil {
		return false
	}
	c := b.makeChatContext(msg.Chat, "Commands::"+command.Name())
	err := command.Run(c, msg, args)
	if err != nil {
		b.log.Error("command %s ended with an error: %s", command.Name(), err.Error())
	}
	return true
}

func (b *Bot) onMessage(msg *tg.Message) {
	if b.cmds != nil {
		cmd, mention, args := ParseCommandMsg(msg)
		if cmd != "" && (mention == "" || mention == b.me.Username) &&
			b.onCommand(msg, cmd, args) {
			return
		}
	}

	h, n := b.handler.OnMessage, "Message"
	if msg.Chat.IsChannel() {
		h, n = b.handler.OnChannelPost, "Post"
	}
	if h != nil {
		h(b.makeChatContext(msg.Chat, n), msg)
	}
}

// Handler conatains all handler functions and handling parameters.
type Handler struct {
	// It only handles messages that are NOT commands.
	OnMessage       func(ChatContext, *tg.Message)
	OnEditedMessage func(ChatContext, *tg.Message)
	// It only handles messages that are NOT commands.
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

func (h *Handler) allowed(b *Bot) []string {
	list := make([]string, 0, 14)
	add := func(a bool, s string) {
		if a {
			list = append(list, s)
		}
	}
	add(h.OnMessage != nil || b.cmds != nil, "message")
	add(h.OnEditedMessage != nil, "edited_message")
	add(h.OnChannelPost != nil || b.cmds != nil, "channel_post")
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
