package commands

import (
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// Filter handles message and if it is a command calls the command handler.
type Filter struct {
	// if empty then any commands with a mention will be
	// considered as not for the current bot.
	Username string

	// if true then command will be handled even it is not for the current bot,
	// otherwise the message will be ignored.
	PassNotForMe bool

	// if true and the PassNotForMe is true then Command will be called
	// for commands not for the current bot, otherwise the message will be passed to Message.
	PassToCommand bool

	// is called when the message is not a command.
	Message func(*tgot.Message, *tg.Message)

	// Command is called when the message is a command for the current bot.
	// The ctx is a child of original context with name 'commands'.
	Command func(m *tgot.Message, msg *tg.Message, cmd string, args []string)
}

// Handle handles message.
func (h *Filter) Handle(ctx *tgot.Message, msg *tg.Message) {
	cmd, mention, args := ParseMsg(msg)
	if cmd == "" || h.Command == nil {
		if h.Message != nil {
			h.Message(ctx, msg)
		}
		return
	}
	if mention != "" && mention != h.Username {
		if !h.PassNotForMe {
			return
		}
		if !h.PassToCommand {
			if h.Message != nil {
				h.Message(ctx, msg)
			}
			return
		}
	}
	h.Command(ctx.WithName("commands"), msg, cmd, args)
}
