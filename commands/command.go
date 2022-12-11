package commands

import (
	"strings"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// Command respresents simple command.
type Command struct {
	Cmd  string
	Func func(tgot.ChatContext, *tg.Message, []string) error

	Description string
	FullDesc    string
	Args        []Arg
}

// Arg type.
type Arg struct {
	Required bool
	Name     string
	Consts   []string
}

// Name returns command name.
func (c Command) Name() string { return c.Cmd }

// Desc returns command description.
func (c Command) Desc() string { return c.Description }

// Run implements tgot.Command and runs command function.
func (c Command) Run(ctx tgot.ChatContext, msg *tg.Message, args []string) error {
	return c.Func(ctx, msg, args)
}

// Help generates help message.
func (c Command) Help() tgot.Message {
	sb := strings.Builder{}
	sb.WriteByte(tgot.Prefix)
	sb.WriteString(c.Cmd)
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
	entities := make([]tg.MessageEntity, 2, 3)
	entities[0] = tg.MessageEntity{
		Type:   tg.EntityCodeBlock,
		Offset: 0,
		Length: sb.Len(),
	}
	sb.WriteString("\n\n")
	sb.WriteString(c.Description)
	entities[1] = tg.MessageEntity{
		Type:   tg.EntityBold,
		Offset: entities[0].Length + 2,
		Length: len(c.Description),
	}
	if len(c.FullDesc) > 0 {
		sb.WriteString("\n\n")
		sb.WriteString(c.FullDesc)
		entities = append(entities, tg.MessageEntity{
			Type:   tg.EntityItalic,
			Offset: sb.Len() - len(c.FullDesc) - 1,
			Length: len(c.FullDesc),
		})
	}

	return tgot.Message{
		Text:     sb.String(),
		Entities: entities,
	}
}

// MakeHelp creates '/help' command.
func MakeHelp(list *List) *Command {
	if list == nil {
		return nil
	}
	h := Command{
		Cmd:         "help",
		Description: "help",
		Args: []Arg{
			{
				Name: "command",
			},
		},
	}
	h.Func = func(ctx tgot.ChatContext, msg *tg.Message, args []string) error {
		if len(args) > 0 {
			if cmd := list.GetCmd(args[0]); cmd != nil {
				return ctx.ReplyE(msg.ID, cmd.Help())
			}
			return ctx.ReplyE(msg.ID, tgot.NewMessage("command not found"))
		}
		var sb strings.Builder
		sb.WriteString("Commands list\n")

		for _, c := range *list {
			sb.WriteByte('\n')
			sb.WriteByte(tgot.Prefix)
			sb.WriteString(c.Cmd + " - " + c.Description)
		}
		return ctx.ReplyE(msg.ID, tgot.NewMessage(sb.String()))
	}
	return &h
}
