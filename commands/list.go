package commands

import (
	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

var _ tgot.Commands = &List{}

// List represents simple commands list.
type List []*Command

// Init sets the default list of the bot's commands on telegram servers.
func (list List) Init(a *api.API) error {
	cmds := make([]tg.Command, len(list))
	for i := range list {
		cmds[i] = tg.Command{
			Command:     list[i].Cmd,
			Description: list[i].Description,
		}
	}
	return a.SetCommands(&api.CommandsData{Commands: cmds})
}

// Get returns command by name.
func (list List) Get(cmd string, _ *tg.Message) tgot.Command {
	c := list.GetCmd(cmd)
	if c == nil {
		return nil
	}
	return c
}

// GetCmd returns command by name.
func (list List) GetCmd(cmd string) *Command {
	for _, c := range list {
		if c.Cmd == cmd {
			return c
		}
	}
	return nil
}
