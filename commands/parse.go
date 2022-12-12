package commands

import (
	"strings"

	"github.com/karalef/tgot/api/tg"
)

// Prefix is the character with which commands must begin.
const Prefix = '/'

// Parse parses and checks the input for command, mention and arguments.
func Parse(c string) (cmd string, mention string, args []string) {
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

// ParseMsg parses and checks the input message for command, mention and arguments.
// In all cases it is faster or equal in performance to Parse.
func ParseMsg(msg *tg.Message) (cmd string, mention string, args []string) {
	text, ents := msg.Text, msg.Entities
	if len(ents) == 0 {
		text, ents = msg.Caption, msg.CaptionEntities
	}
	if len(ents) == 0 || ents[0].Type != tg.EntityCommand || ents[0].Offset != 0 {
		return "", "", nil
	}

	cmd = text[:ents[0].Length]

	if len(text) > len(cmd)+1 {
		args = strings.Split(text[len(cmd)+1:], " ")
	}

	if i := strings.Index(cmd, "@"); i != -1 && len(cmd) > i+1 {
		mention = cmd[i+1:]
		cmd = cmd[:i]
	}
	cmd = cmd[1:]
	return
}
