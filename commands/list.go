package commands

import (
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// List represents simple commands list.
type List []*Command

// Setup sets the default list of the bot's commands on telegram servers.
func (list *List) Setup(b *tgot.Bot) error {
	cmds := make([]tg.Command, len(*list))
	for i := range *list {
		cmds[i] = tg.Command{
			Command:     (*list)[i].Cmd,
			Description: (*list)[i].Description,
		}
	}
	return b.SetCommands(&tgot.CommandsData{Commands: cmds})
}

// Command runs a command if it exists.
func (list *List) Command(ctx tgot.ChatContext, msg *tg.Message, cmd string, args []string) error {
	c := list.GetCmd(cmd)
	if c == nil {
		return nil
	}
	return c.Run(ctx.Child(cmd), msg, args)
}

// GetCmd returns command by name.
func (list *List) GetCmd(cmd string) *Command {
	for _, c := range *list {
		if c.Cmd == cmd {
			return c
		}
	}
	return nil
}
