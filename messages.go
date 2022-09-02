package bot

import (
	"tghwbot/bot/tg"
)

func (b *Bot) makeMessageContext(msg *tg.Message, name string) MessageContext {
	return MessageContext{
		ChatContext: b.makeChatContext(msg.Chat, name),
		msgID:       msg.ID,
	}
}

// MessageContext type.
type MessageContext struct {
	ChatContext
	msgID int
}

// Reply replies to a message with the Sendable and returns only an error.
func (c MessageContext) Reply(s Sendable) error {
	_, err := c.Chat.Send(s, SendOptions[tg.ReplyMarkup]{
		BaseSendOptions: BaseSendOptions{ReplyTo: c.msgID},
	})
	return err
}

// ReplyText replies to a message with just a text and returns only an error.
func (c MessageContext) ReplyText(text string) error {
	return c.Reply(NewMessage(text))
}

// MessageSignature creates a chat message signature.
//
// Returns nil if the message is not sent by the current bot.
func (c Context) MessageSignature(msg *tg.Message) *MessageSignature {
	if msg.From.ID != c.bot.me.ID {
		return nil
	}
	return &MessageSignature{
		chatID: msg.Chat.ID,
		msgID:  msg.ID,
	}
}

// InlineSignature creates a chat message signature.
func (c Context) InlineSignature(i *tg.InlineChosen) *MessageSignature {
	return &MessageSignature{inline: i.InlineMessageID}
}

func (c Context) sig(sig *MessageSignature, method string, p params, files ...file) (*tg.Message, error) {
	sig.signature(p)
	if sig.isInline() {
		return nil, c.api(method, p, files...)
	}
	return api[*tg.Message](c, method, p, files...)
}

// MessageSignature contains inline message id or chat id with message id.
type MessageSignature struct {
	chatID int64
	msgID  int
	inline string
}

func (s MessageSignature) isInline() bool {
	return s.inline != ""
}

func (s MessageSignature) signature(p params) {
	if s.inline != "" {
		p.set("inline_message_id", s.inline)
		return
	}
	p.setInt64("chat_id", s.chatID)
	p.setInt("message_id", s.msgID)
}

// LiveLocation contains parameters for editing the live location.
type LiveLocation struct {
	Long               float32
	Lat                float32
	HorizontalAccuracy *float32
	Heading            int
	AlertRadius        int
}

// EditLiveLocation edits live location messages.
func (c Context) EditLiveLocation(sig *MessageSignature, l LiveLocation, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	p := params{}
	p.setFloat("latitude", l.Lat)
	p.setFloat("longitude", l.Long)
	if l.HorizontalAccuracy != nil {
		p.setFloat("horizontal_accuracy", *l.HorizontalAccuracy)
	}
	p.setInt("heading", l.Heading)
	p.setInt("proximity_alert_radius", l.AlertRadius)
	if len(replyMarkup) > 0 {
		p.setJSON("reply_markup", replyMarkup[0])
	}
	return c.sig(sig, "editMessageLiveLocation", p)
}

// StopLiveLocation stops updating a live location message before live_period expires.
func (c Context) StopLiveLocation(sig *MessageSignature, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	p := params{}
	if len(replyMarkup) > 0 {
		p.setJSON("reply_markup", replyMarkup[0])
	}
	return c.sig(sig, "stopMessageLiveLocation", p)
}

// EditText contains parameters for editing message text.
type EditText struct {
	Text                  string
	ParseMode             tg.ParseMode
	Entities              []tg.MessageEntity
	DisableWebPagePreview bool
}

// EditText edits text and game messages.
func (c Context) EditText(sig *MessageSignature, t EditText, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	p := params{}
	p.set("text", t.Text)
	p.set("parse_mode", string(t.ParseMode))
	p.setJSON("entities", t.Entities)
	p.setBool("disable_web_page_preview", t.DisableWebPagePreview)
	if len(replyMarkup) > 0 {
		p.setJSON("reply_markup", replyMarkup[0])
	}
	return c.sig(sig, "editMessageText", p)
}

// EditCaption edits captions of messages.
func (c Context) EditCaption(sig *MessageSignature, cap CaptionData, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	p := params{}
	cap.embed(p)
	if len(replyMarkup) > 0 {
		p.setJSON("reply_markup", replyMarkup[0])
	}
	return c.sig(sig, "editMessageCaption", p)
}

// EditMedia edits animation, audio, document, photo, or video messages.
func (c Context) EditMedia(sig *MessageSignature, m tg.MediaInputter, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	p := params{}
	files, err := prepareInputMedia(p, false, m)
	if err != nil {
		return nil, err
	}
	if len(replyMarkup) > 0 {
		p.setJSON("reply_markup", replyMarkup[0])
	}
	return c.sig(sig, "editMessageMedia", p, files...)
}

// EditReplyMarkup edits only the reply markup of messages.
func (c Context) EditReplyMarkup(sig *MessageSignature, replyMarkup *tg.InlineKeyboardMarkup) (*tg.Message, error) {
	p := params{}
	p.setJSON("reply_markup", replyMarkup)
	return c.sig(sig, "editMessageReplyMarkup", p)
}
