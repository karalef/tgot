package bot

import (
	"strings"
	"tghwbot/bot/tg"
)

func makeHelp(b *Bot) *Command {
	h := Command{
		Cmd:         "help",
		Description: "help",
		Args: []Arg{
			{
				Name: "command",
			},
		},
	}
	h.Run = func(ctx MessageContext, msg *tg.Message, args []string) error {
		if len(args) > 0 {
			for _, c := range b.cmds {
				if c.Cmd != args[0] {
					continue
				}
				return ctx.Reply(c.GenerateHelp())
			}
			return ctx.ReplyText("command not found")
		}
		var sb strings.Builder
		sb.WriteString("Commands list\n")

		for _, c := range b.cmds {
			sb.WriteByte('\n')
			sb.WriteByte(Prefix)
			sb.WriteString(c.Cmd + " - " + c.Description)
		}
		return ctx.ReplyText(sb.String())
	}
	return &h
}
