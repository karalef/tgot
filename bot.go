package tgot

import (
	"errors"
	"net/http"
	"sync"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/internal"
	"github.com/karalef/tgot/logger"
	"github.com/karalef/tgot/tg"
	"github.com/karalef/tgot/updates"
)

// Config contains bot configuration.
type Config struct {
	APIURL   string         // telegram default if empty
	FileURL  string         // telegram default if empty
	Client   *http.Client   // http.DefaultClient if empty
	Logger   *logger.Logger // logger.Default if empty
	Handler  Handler
	Commands []*Command
	Options  Options
}

// Options contains bot options.
type Options = internal.Flag

// available options.
const (
	CaptionCommands Options = 1 << iota
	ChannelCommands
	MakeHelp
)

// New creates new bot.
func New(token string, poller updates.Poller, config Config) (*Bot, error) {
	if poller == nil {
		return nil, errors.New("nil poller")
	}
	a, me, err := api.New(token, config.APIURL, config.FileURL, config.Client)
	if err != nil {
		return nil, err
	}
	if config.Logger == nil {
		config.Logger = logger.Default("TGO")
	}
	b := Bot{
		api:     a,
		log:     config.Logger,
		poller:  poller,
		opts:    config.Options,
		handler: config.Handler,
		cmds:    config.Commands,
		me:      *me,
	}
	if config.Options.Has(MakeHelp) {
		b.cmds = append(b.cmds, makeHelp(&b))
	}

	return &b, nil
}

// Bot type.
type Bot struct {
	api    *api.API
	log    *logger.Logger
	poller updates.Poller
	opts   Options

	cls sync.Once

	handler Handler
	cmds    []*Command

	me tg.User
}

func (b *Bot) setupCommands() error {
	commands := make([]tg.Command, len(b.cmds))
	for i := range b.cmds {
		commands[i] = tg.Command{
			Command:     b.cmds[i].Cmd,
			Description: b.cmds[i].Description,
		}
	}
	return b.api.SetCommands(&api.CommandsData{Commands: commands})
}

func (b *Bot) cancel(err error) {
	b.cls.Do(func() {
		if err != nil {
			b.poller.Close()
			b.log.Error(err.Error())
		} else {
			b.Stop()
		}
	})
}

// Stop stops polling for updates.
func (b *Bot) Stop() {
	b.poller.Shutdown()
}

// Run starts bot.
func (b *Bot) Run() error {
	err := b.setupCommands()
	if err != nil {
		return err
	}

	allowed := b.handler.allowed(b)
	return b.poller.Run(b.api, b.handle, allowed)
}

func (b *Bot) handle(upd *tg.Update) {
	h := &b.handler

	// since [Handler.allowed] returns a list of the specified handlers
	// it makes no sense to check for nil except for OnMessage and OnChannelPost.
	switch {
	case upd.Message != nil:
		b.onMessage(upd.Message)
	case upd.EditedMessage != nil:
		ctx := b.makeMessageContext(upd.EditedMessage, "EditedMessage")
		h.OnEditedMessage(ctx, upd.EditedMessage)
	case upd.ChannelPost != nil:
		b.onChannelPost(upd.ChannelPost)
	case upd.EditedChannelPost != nil:
		ctx := b.makeMessageContext(upd.EditedChannelPost, "EditedPost")
		h.OnEditedChannelPost(ctx, upd.EditedChannelPost)
	case upd.CallbackQuery != nil:
		h.OnCallbackQuery(CallbackContext{
			Context: b.MakeContext("Callback"),
			queryID: upd.CallbackQuery.ID,
		}, upd.CallbackQuery)
	case upd.InlineQuery != nil:
		h.OnInlineQuery(InlineContext{
			Context: b.MakeContext("Inline"),
			queryID: upd.InlineQuery.ID,
		}, upd.InlineQuery)
	case upd.InlineChosen != nil:
		h.OnInlineChosen(b.MakeContext("InlineChosen"), upd.InlineChosen)
	case upd.ShippingQuery != nil:
		h.OnShippingQuery(ShippingContext{
			Context: b.MakeContext("ShippingQuery"),
			queryID: upd.ShippingQuery.ID,
		}, upd.ShippingQuery)
	case upd.PreCheckoutQuery != nil:
		h.OnPreCheckoutQuery(PreCheckoutContext{
			Context: b.MakeContext("PreCheckoutQuery"),
			queryID: upd.PreCheckoutQuery.ID,
		}, upd.PreCheckoutQuery)
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

func (b *Bot) parseCommand(c string) (cmd *Command, args []string) {
	c, mention, args := ParseCommand(c)
	if mention != "" && mention != b.me.Username {
		return nil, nil
	}
	for _, cmd := range b.cmds {
		if c == cmd.Cmd {
			return cmd, args
		}
	}
	return nil, nil
}

func (b *Bot) parseMessage(msg *tg.Message) (cmd *Command, args []string) {
	text := msg.Text
	if text == "" {
		if !b.opts.Has(CaptionCommands) {
			return
		}
		text = msg.Caption
	}
	return b.parseCommand(text)
}

func (b *Bot) onMessage(msg *tg.Message) {
	cmd, args := b.parseMessage(msg)
	if cmd != nil {
		b.onCommand(msg, cmd, args)
		return
	}

	if b.handler.OnMessage != nil {
		b.handler.OnMessage(b.makeMessageContext(msg, "Message"), msg)
	}
}

func (b *Bot) onChannelPost(msg *tg.Message) {
	if b.opts.Has(ChannelCommands) {
		cmd, args := b.parseMessage(msg)
		if cmd != nil {
			b.onCommand(msg, cmd, args)
			return
		}
	}

	if b.handler.OnChannelPost != nil {
		b.handler.OnChannelPost(b.makeMessageContext(msg, "Post"), msg)
	}
}

// Handler conatains all handler functions and handling parameters.
type Handler struct {
	// It only handles messages that are NOT commands.
	OnMessage       func(MessageContext, *tg.Message)
	OnEditedMessage func(MessageContext, *tg.Message)
	// It only handles messages that are NOT commands
	// or if commands in the channels are disabled.
	OnChannelPost       func(MessageContext, *tg.Message)
	OnEditedChannelPost func(MessageContext, *tg.Message)
	OnCallbackQuery     func(CallbackContext, *tg.CallbackQuery)
	OnInlineQuery       func(InlineContext, *tg.InlineQuery)
	OnInlineChosen      func(Context, *tg.InlineChosen)
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
	add(h.OnMessage != nil || len(b.cmds) > 0, "message")
	add(h.OnEditedMessage != nil, "edited_message")
	add(h.OnChannelPost != nil || (len(b.cmds) > 0 && b.opts.Has(ChannelCommands)), "channel_post")
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
