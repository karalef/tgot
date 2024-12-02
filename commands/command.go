package commands

import (
	"strings"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// Command represents command.
type Command interface {
	Name() string
	Description() string
	Run(*tgot.Message, *tg.Message, []string)

	// Is returns true if this command matches the given string.
	Is(string) bool

	// Help generates help message.
	Help() tgot.Text
}

var _ Command = SimpleCommand{}

// SimpleCommand represents simple command.
type SimpleCommand struct {
	Command string
	Aliases []string
	Func    func(*tgot.Message, *tg.Message, []string) error

	Args     []Arg
	Desc     string
	FullDesc string
}

// Arg type.
type Arg struct {
	Name     string
	Consts   []string
	Required bool
}

// Name returns command name.
func (c SimpleCommand) Name() string { return c.Command }

// Description returns command description.
func (c SimpleCommand) Description() string { return c.Desc }

// Run runs command.
func (c SimpleCommand) Run(m *tgot.Message, msg *tg.Message, args []string) {
	if c.Func != nil {
		c.Func(m, msg, args)
	}
}

// Is returns true if this command matches the given string.
func (c SimpleCommand) Is(cmd string) bool {
	if c.Command == cmd {
		return true
	}
	for _, a := range c.Aliases {
		if a == cmd {
			return true
		}
	}
	return false
}

// Help generates help message.
func (c SimpleCommand) Help() tgot.Text {
	sb := strings.Builder{}
	entities := make([]tg.MessageEntity, 2, 3)

	// description
	sb.WriteString(c.Command + " - " + c.Desc)
	entities[0] = tg.MessageEntity{
		Type:   tg.EntityBold,
		Offset: 0,
		Length: sb.Len(),
	}

	// usage
	sb.WriteString("\n\nUsage:\n")
	entities[1] = tg.MessageEntity{
		Type:   tg.EntityCodeBlock,
		Offset: sb.Len(),
	}
	sb.WriteByte(Prefix)
	sb.WriteString(c.Command)
	for _, a := range c.Args {
		sb.WriteByte(' ')
		if a.Required {
			sb.WriteByte('[')
		} else {
			sb.WriteByte('{')
		}
		sb.WriteString(a.Name)

		if len(a.Consts) > 0 {
			if len(a.Name) > 0 {
				sb.WriteByte(':')
			}
			sb.WriteString("\"" + strings.Join(a.Consts, "\"|\"") + "\"")
		}
		if a.Required {
			sb.WriteByte(']')
		} else {
			sb.WriteByte('}')
		}
	}
	entities[1].Length = sb.Len() - entities[1].Offset

	// full description
	if len(c.FullDesc) > 0 {
		sb.WriteString("\n\n")
		entities = append(entities, tg.MessageEntity{
			Type:   tg.EntityItalic,
			Offset: sb.Len(),
			Length: len(c.FullDesc),
		})
		sb.WriteString(c.FullDesc)
	}

	return tgot.Text{
		Text:     sb.String(),
		Entities: entities,
	}
}

// MakeHelp creates '/help' command.
func MakeHelp(list *List) SimpleCommand {
	h := SimpleCommand{
		Command: "help",
		Desc:    "help",
		Args: []Arg{
			{
				Name: "command",
			},
		},
	}
	h.Func = func(m *tgot.Message, msg *tg.Message, args []string) error {
		chat := tgot.WithChatID(m, m.ID().ChatID())
		if len(args) > 0 {
			if cmd := list.GetCmd(args[0]); cmd != nil {
				return chat.ReplyE(tgot.ReplyTo(msg.ID), cmd.Help())
			}
			return chat.ReplyE(tgot.ReplyTo(msg.ID), tgot.NewText("command not found"))
		}
		var sb strings.Builder
		sb.WriteString("Commands list\n")

		for _, c := range *list {
			sb.WriteByte('\n')
			sb.WriteByte(Prefix)
			sb.WriteString(c.Name() + " - " + c.Description())
		}
		return chat.ReplyE(tgot.ReplyTo(msg.ID), tgot.NewText(sb.String()))
	}
	return h
}
