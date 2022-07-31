package bot

import "tghwbot/bot/tg"

// Chat represents chat api.
type Chat struct {
	ctx    commonContext
	chatID string
}

// GetInfo returns up to date information about the chat.
func (c *Chat) GetInfo() *tg.Chat {
	return api[*tg.Chat](c.ctx, "getChat", params{
		"chat_id": {c.chatID},
	})
}

// MemberCount returns the number of members in a chat.
func (c *Chat) MemberCount() int {
	return api[int](c.ctx, "getChatMemberCount", params{
		"chat_id": {c.chatID},
	})
}

// Leave a group, supergroup or channel.
func (c *Chat) Leave() {
	api[bool](c.ctx, "leaveChat", params{
		"chat_id": {c.chatID},
	})
}

// Forward forwards messages of any kind.
// Service messages can't be forwarded.
func (c *Chat) Forward(from *Chat, msgID int) *tg.Message {
	p := params{}.set("chat_id", c.chatID)
	p.set("from_chat_id", from.chatID)
	p.set("message_id", msgID)
	p.set("disable_notification", false)
	p.set("protect_content", false)
	return api[*tg.Message](c.ctx, "forwardMessage", p)
}

// ForwardTo forwards message to specified chat instead of current.
func (c *Chat) ForwardTo(to *Chat, msgID int) *tg.Message {
	return to.Forward(c, msgID)
}

// SendText sends just a text.
func (c *Chat) SendText(text string) *tg.Message {
	return c.Send(NewMessage(text))
}

// Send sends any Sendable object.
func (c *Chat) Send(s Sendable, opts ...SendOptions) *tg.Message {
	if s == nil {
		return nil
	}

	m := "send" + s.what()
	p := s.params().set("chat_id", c.chatID)
	if len(opts) > 0 {
		o := &opts[0]
		p.set("disable_notification", o.DisableNotification)
		p.set("protect_content", o.ProtectContent)
		p.set("reply_to_message_id", o.ReplyTo)
		p.set("allow_sending_without_reply", o.AllowSendingWithoutReply)
		p.set("reply_markup", o.ReplyMarkup)
	}
	if f, ok := s.(Fileable); ok {
		files := f.files()
		for i := range files {
			if files[i].Data() == "" {
				return api[*tg.Message](c.ctx, m, p, files...)
			}
		}
		for i := range files {
			p.set(files[i].Field, files[i].Data())
		}
	}

	return api[*tg.Message](c.ctx, m, p)
}

// SendChatAction sends chat action to tell the user that something
// is happening on the bot's side.
func (c *Chat) SendChatAction(act tg.ChatAction) {
	p := params{}
	p.set("chat_id", c.chatID)
	p.set("action", string(act))
	api[bool](c.ctx, "sendChatAction", p)
}

// SendInvoice sends invoice.
func (c *Chat) SendInvoice(i Invoice) *tg.Message {
	p := i.params().set("chat_id", c.chatID)
	return api[*tg.Message](c.ctx, "sendInvoice", p)
}
