package bot

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"tghwbot/bot/logger"
	"tghwbot/bot/tg"
)

// Config contains bot configuration.
type Config struct {
	APIURL   string         // tg.DefaultAPIURL if empty
	FileURL  string         // tg.DefaultFileURL if empty
	Client   *http.Client   // http.DefaultClient if empty
	Logger   *logger.Logger // logger.Default if empty
	Handler  Handler
	Commands []*Command
	MakeHelp bool
}

// New creates new bot.
func New(token string, config Config) (*Bot, error) {
	if token == "" {
		return nil, errors.New("no token provided")
	}
	if config.APIURL == "" {
		config.APIURL = tg.DefaultAPIURL
	}
	if config.FileURL == "" {
		config.FileURL = tg.DefaultFileURL
	}
	if config.Client == nil {
		config.Client = http.DefaultClient
	}
	if config.Logger == nil {
		config.Logger = logger.Default("TGO")
	}
	b := Bot{
		token:   token,
		apiURL:  config.APIURL,
		fileURL: config.FileURL,
		client:  config.Client,
		log:     config.Logger,
		cmds:    config.Commands,
	}
	if config.MakeHelp {
		b.cmds = append(b.cmds, makeHelp(&b))
	}

	me, err := b.GetMe()
	if err != nil {
		return nil, err
	}
	b.me = *me

	return &b, nil
}

// Bot type.
type Bot struct {
	token   string
	apiURL  string
	fileURL string
	client  *http.Client
	log     *logger.Logger

	wg      sync.WaitGroup
	stop    context.CancelFunc
	err     error
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
	return b.setCommands(&commandParams{Commands: commands})
}

func (b *Bot) cancel(err error) {
	if b.err == nil {
		b.err = err
	}
	b.stop()
}

// Stop stops polling for updates.
// The call is similar to context cancellation.
// It panics if the bot has not been run.
func (b *Bot) Stop() {
	b.cancel(nil)
	b.wg.Wait()
}

// Run starts bot.
func (b *Bot) Run() error {
	return b.RunContext(context.Background())
}

// RunContext starts bot.
// Returns nil if context is closed.
func (b *Bot) RunContext(ctx context.Context) error {
	if b.stop != nil {
		return errors.New("bot is already running")
	}
	ctx, b.stop = context.WithCancel(ctx)

	err := b.setupCommands()
	if err != nil {
		return err
	}

	allowed := b.handler.allowed(b)
	defer b.wg.Wait()
	offset := 0
	for {
		upds, err := b.getUpdates(ctx, offset+1, 30, 0, allowed)
		switch err {
		case nil:
		case context.Canceled, context.DeadlineExceeded:
			return nil
		default:
			return err
		}
		for i := range upds {
			go b.handle(&upds[i])
		}
		offset = upds[len(upds)-1].ID
	}
}

func (b *Bot) handle(upd *tg.Update) {
	b.wg.Add(1)
	defer b.wg.Done()
	h := &b.handler
	if h.Filter != nil && !h.Filter(upd) {
		return
	}

	// since [Handler.allowed] returns a list of the specified handlers
	// it makes no sense to check for nil except for OnMessage and OnChannelPost.
	switch {
	case upd.Message != nil:
		b.onMessage(upd.Message)
	case upd.EditedMessage != nil:
		ctx := b.makeMessageContext(upd.EditedMessage, "Edited message handler")
		h.OnEditedMessage(ctx, upd.EditedMessage)
	case upd.ChannelPost != nil:
		b.onChannelPost(upd.ChannelPost)
	case upd.EditedChannelPost != nil:
		ctx := b.makeMessageContext(upd.EditedChannelPost, "Edited channel post handler")
		h.OnEditedChannelPost(ctx, upd.EditedChannelPost)
	case upd.CallbackQuery != nil:
		h.OnCallbackQuery(QueryContext[CallbackAnswer]{
			Context: b.makeContext("Callback handler"),
			queryID: upd.CallbackQuery.ID,
		}, upd.CallbackQuery)
	case upd.InlineQuery != nil:
		h.OnInlineQuery(QueryContext[InlineAnswer]{
			Context: b.makeContext("Inline handler"),
			queryID: upd.InlineQuery.ID,
		}, upd.InlineQuery)
	case upd.InlineChosen != nil:
		h.OnInlineChosen(b.makeContext("Inline chosen handler"), upd.InlineChosen)
	case upd.ShippingQuery != nil:
		h.OnShippingQuery(QueryContext[ShippingAnswer]{
			Context: b.makeContext("Shipping query handler"),
			queryID: upd.ShippingQuery.ID,
		}, upd.ShippingQuery)
	case upd.PreCheckoutQuery != nil:
		h.OnPreCheckoutQuery(QueryContext[PreCheckoutAnswer]{
			Context: b.makeContext("Pre checkout query handler"),
			queryID: upd.PreCheckoutQuery.ID,
		}, upd.PreCheckoutQuery)
	case upd.Poll != nil:
		h.OnPoll(b.makeContext("Poll handler"), upd.Poll)
	case upd.PollAnswer != nil:
		h.OnPollAnswer(b.makeContext("Poll answer handler"), upd.PollAnswer)
	case upd.MyChatMember != nil:
		ctx := b.makeChatContext(upd.MyChatMember.Chat, "My chat member handler")
		h.OnMyChatMember(ctx, upd.MyChatMember)
	case upd.ChatMember != nil:
		ctx := b.makeChatContext(upd.ChatMember.Chat, "Chat member handler")
		h.OnChatMember(ctx, upd.ChatMember)
	case upd.ChatJoinRequest != nil:
		ctx := b.makeChatContext(upd.ChatJoinRequest.Chat, "Chat join request handler")
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
		if b.handler.DisableCaptionCommands {
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
		b.handler.OnMessage(b.makeMessageContext(msg, "Message handler"), msg)
	}
}

func (b *Bot) onChannelPost(msg *tg.Message) {
	if b.handler.ChannelCommands {
		cmd, args := b.parseMessage(msg)
		if cmd != nil {
			b.onCommand(msg, cmd, args)
			return
		}
	}

	if b.handler.OnChannelPost != nil {
		b.handler.OnChannelPost(b.makeMessageContext(msg, "Channel post handler"), msg)
	}
}

// Handler conatains all handler functions and handling parameters.
type Handler struct {
	DisableCaptionCommands bool
	ChannelCommands        bool

	// It should return false if the update is not to be handled.
	Filter func(*tg.Update) (pass bool)

	// It only handles messages that are NOT commands.
	OnMessage       func(MessageContext, *tg.Message)
	OnEditedMessage func(MessageContext, *tg.Message)
	// It only handles messages that are NOT commands
	// or if commands in the channels are disabled by [Handler.ChannelCommands].
	OnChannelPost       func(MessageContext, *tg.Message)
	OnEditedChannelPost func(MessageContext, *tg.Message)
	OnCallbackQuery     func(QueryContext[CallbackAnswer], *tg.CallbackQuery)
	OnInlineQuery       func(QueryContext[InlineAnswer], *tg.InlineQuery)
	OnInlineChosen      func(Context, *tg.InlineChosen)
	OnShippingQuery     func(QueryContext[ShippingAnswer], *tg.ShippingQuery)
	OnPreCheckoutQuery  func(QueryContext[PreCheckoutAnswer], *tg.PreCheckoutQuery)
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
	add(h.OnChannelPost != nil || (len(b.cmds) > 0 && h.ChannelCommands), "channel_post")
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
