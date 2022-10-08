package tgot

import (
	"strings"

	"github.com/karalef/tgot/tg"
)

// Prefix is the character with which commands must begin.
const Prefix = '/'

func (b *Bot) onCommand(msg *tg.Message, cmd *Command, args []string) {
	c := b.makeMessageContext(msg, "Commands::"+cmd.Cmd)
	err := cmd.Run(c, msg, args)
	if err != nil {
		c.bot.log.Error("command %s ended with an error: %s", cmd.Cmd, err.Error())
	}
}

// Command respresents conversation command.
type Command struct {
	Cmd string
	Run func(MessageContext, *tg.Message, []string) error

	Description string
	Help        string
	Args        []Arg
}

// Arg type.
type Arg struct {
	Required bool
	Name     string
	Consts   []string
}

// GenerateHelp generates help message.
func (c *Command) GenerateHelp() Message {
	sb := strings.Builder{}
	sb.WriteByte(Prefix)
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
	if len(c.Help) > 0 {
		sb.WriteString("\n\n")
		sb.WriteString(c.Help)
		entities = append(entities, tg.MessageEntity{
			Type:   tg.EntityItalic,
			Offset: sb.Len() - len(c.Help) - 1,
			Length: len(c.Help),
		})
	}

	return Message{
		Text:     sb.String(),
		Entities: entities,
	}
}

// ParseCommand parses and checks the input for command, mention and arguments.
func ParseCommand(c string) (cmd string, mention string, args []string) {
	if len(c) < 2 || c[0] != Prefix {
		return "", "", nil
	}
	split := strings.Split(c[1:], " ")
	cmd, args = split[0], split[1:]
	if i := strings.Index(cmd, "@"); i != -1 && len(cmd) > i+1 {
		mention = cmd[i+1:]
		cmd = cmd[:i]
	}
	return
}

// ParseCommandMsg parses and checks the input message for command, mention and arguments.
// In all cases it is faster or equal in performance to ParseCommand.
func ParseCommandMsg(msg *tg.Message) (cmd string, mention string, args []string) {
	ents := msg.Entities
	if len(ents) == 0 || ents[0].Type != tg.EntityCommand || ents[0].Offset != 0 {
		return "", "", nil
	}
	cmd = msg.Text[:ents[0].Length]
	if len(msg.Text) > len(cmd)+1 {
		args = strings.Split(msg.Text[ents[0].Length+1:], " ")
	}
	if i := strings.Index(cmd, "@"); i != -1 && len(cmd) > i+1 {
		mention = cmd[i+1:]
		cmd = cmd[:i]
	}
	return
}
