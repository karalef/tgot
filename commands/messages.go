package commands

import (
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// MessageHandler handles message and if it is a command calls the command handler.
type MessageHandler struct {
	// if empty then any commands with a mention will be
	// considered as not for the current bot
	Username string

	// if true then Message will be called for commands not for the current bot,
	// otherwise the message will be ignored
	PassNotForMe bool

	// is called when the message is not a command
	Message func(tgot.ChatContext, *tg.Message)

	// Command is called when the message is a command for the current bot.
	// It creates a child context with name 'Commands' and logs an error.
	Command func(ctx tgot.ChatContext, msg *tg.Message, cmd string, args []string) error
}

// Handle handles message.
func (h *MessageHandler) Handle(ctx tgot.ChatContext, msg *tg.Message) {
	cmd, mention, args := ParseMsg(msg)
	if cmd == "" || h.Command == nil {
		if h.Message != nil {
			h.Message(ctx, msg)
		}
		return
	}
	if mention != "" && mention != h.Username {
		if h.PassNotForMe && h.Message != nil {
			h.Message(ctx, msg)
		}
		return
	}
	ctx = ctx.Child("Commands")
	err := h.Command(ctx, msg, cmd, args)
	if err != nil {
		ctx.Logger().Error("command '%s' ended with an error: %s", cmd, err.Error())
	}
}
