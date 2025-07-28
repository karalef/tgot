package cmd

import (
	"strings"

	"github.com/karalef/tgot/api/tg"
)

// Command represents parsed command.
type Command struct {
	Name    string
	Mention string
	Args    []string
}

// Prefix is the character with which commands must begin.
const Prefix = '/'

// Parse parses and checks the input for command, mention and arguments.
func Parse(c string) (cmd Command) {
	if len(c) < 2 || c[0] != Prefix {
		return cmd
	}
	split := strings.Split(c[1:], " ")
	cmd.Name, cmd.Args = split[0], split[1:]
	if i := strings.Index(cmd.Name, "@"); i != -1 && len(cmd.Name) > i+1 {
		cmd.Mention = cmd.Name[i+1:]
		cmd.Name = cmd.Name[:i]
	}
	return
}

// ParseMsg parses and checks the input message for command, mention and
// arguments.
// It is wrapper for ParseEntities and automatically detects text and caption.
func ParseMsg(msg *tg.Message) (cmd Command) {
	text, ents := msg.Text, msg.Entities
	if len(ents) == 0 {
		text, ents = msg.Caption, msg.CaptionEntities
	}

	return ParseEntities(text, ents)
}

// ParseEntities parses and checks the input text and entities for command,
// mention and arguments.
// It uses API's entities instead of indexing split symbols.
func ParseEntities(text string, ents []tg.MessageEntity) (cmd Command) {
	if len(ents) == 0 || ents[0].Type != tg.EntityCommand || ents[0].Offset != 0 {
		return cmd
	}

	cmd.Name = text[:ents[0].Length]

	if len(text) > len(cmd.Name)+1 {
		cmd.Args = strings.Split(text[len(cmd.Name)+1:], " ")
	}

	if i := strings.Index(cmd.Name, "@"); i != -1 && len(cmd.Name) > i+1 {
		cmd.Mention = cmd.Name[i+1:]
		cmd.Name = cmd.Name[:i]
	}
	cmd.Name = cmd.Name[1:]
	return
}
