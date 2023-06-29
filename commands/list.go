package commands

import (
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// List represents simple commands list.
type List []Command

// Setup sets the default list of the bot's commands on Telegram servers.
func (list List) Setup(b *tgot.Bot) error {
	cmds := make([]tg.Command, len(list))
	for i := range list {
		cmds[i] = tg.Command{
			Command:     list[i].Name(),
			Description: list[i].Description(),
		}
	}
	return b.SetCommands(&tgot.CommandsData{Commands: cmds})
}

// Command runs a command if it exists.
func (list List) Command(ctx tgot.ChatContext, msg *tg.Message, cmd string, args []string) {
	c := list.GetCmd(cmd)
	if c != nil {
		c.Run(ctx.Child(c.Name()), msg, args)
	}
}

// Exists checks if command exists.
func (list List) Exists(cmd string) bool {
	return list.GetCmd(cmd) != nil
}

// GetCmd returns command by name.
func (list List) GetCmd(cmd string) Command {
	for _, c := range list {
		if c.Is(cmd) {
			return c
		}
	}
	return nil
}
