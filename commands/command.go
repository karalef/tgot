package commands

import (
	"strings"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// Command represents simple command.
type Command struct {
	Name    string
	Aliases []string
	Func    func(tgot.ChatContext, *tg.Message, []string) error

	Args            []Arg
	Description     string
	FullDescription string
}

// Arg type.
type Arg struct {
	Required bool
	Name     string
	Consts   []string
}

// Help generates help message.
func (c Command) Help() tgot.Message {
	sb := strings.Builder{}
	entities := make([]tg.MessageEntity, 2, 3)

	// description
	sb.WriteString(c.Name + " - " + c.Description)
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
	sb.WriteString(c.Name)
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
	if len(c.FullDescription) > 0 {
		sb.WriteString("\n\n")
		entities = append(entities, tg.MessageEntity{
			Type:   tg.EntityItalic,
			Offset: sb.Len(),
			Length: len(c.FullDescription),
		})
		sb.WriteString(c.FullDescription)
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
		Name:        "help",
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
			sb.WriteByte(Prefix)
			sb.WriteString(c.Name + " - " + c.Description)
		}
		return ctx.ReplyE(msg.ID, tgot.NewMessage(sb.String()))
	}
	return &h
}
