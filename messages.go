package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/tg"
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

// MsgSignature creates a chat message signature.
func MsgSignature(msg *tg.Message) MessageSignature {
	return MessageSignature{
		ChatID: msg.Chat.ID,
		MsgID:  msg.ID,
	}
}

// InlineSignature creates an inline message signature.
func InlineSignature(i *tg.InlineChosen) MessageSignature {
	return MessageSignature{Inline: i.InlineMessageID}
}

// CallbackSignature creates a MessageSignature of any callback message type.
func CallbackSignature(q *tg.CallbackQuery) MessageSignature {
	if q.Message != nil {
		return MsgSignature(q.Message)
	}
	return MessageSignature{Inline: q.InlineMessageID}
}

func (c Context) sig(sig MessageSignature, meth string, d api.Data) (*tg.Message, error) {
	sig.signature(&d)
	if sig.isInline() {
		return nil, c.method(meth, d)
	}
	return method[*tg.Message](c, meth, d)
}

// MessageSignature contains inline message id or chat id with message id.
type MessageSignature struct {
	ChatID int64
	MsgID  int
	Inline string
}

func (s MessageSignature) isInline() bool {
	return s.Inline != ""
}

func (s MessageSignature) signature(d *api.Data) {
	if s.Inline != "" {
		d.Set("inline_message_id", s.Inline)
		return
	}
	d.SetInt64("chat_id", s.ChatID)
	d.SetInt("message_id", s.MsgID)
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
func (c Context) EditLiveLocation(sig MessageSignature, l LiveLocation, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	d.SetFloat("latitude", l.Lat)
	d.SetFloat("longitude", l.Long)
	if l.HorizontalAccuracy != nil {
		d.SetFloat("horizontal_accuracy", *l.HorizontalAccuracy, true)
	}
	d.SetInt("heading", l.Heading)
	d.SetInt("proximity_alert_radius", l.AlertRadius)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.sig(sig, "editMessageLiveLocation", d)
}

// StopLiveLocation stops updating a live location message before live_period expires.
func (c Context) StopLiveLocation(sig MessageSignature, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.sig(sig, "stopMessageLiveLocation", d)
}

// EditText contains parameters for editing message text.
type EditText struct {
	Text                  string
	ParseMode             tg.ParseMode
	Entities              []tg.MessageEntity
	DisableWebPagePreview bool
}

// EditText edits text and game messages.
func (c Context) EditText(sig MessageSignature, t EditText, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	d.Set("text", t.Text)
	d.Set("parse_mode", string(t.ParseMode))
	d.SetJSON("entities", t.Entities)
	d.SetBool("disable_web_page_preview", t.DisableWebPagePreview)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.sig(sig, "editMessageText", d)
}

// EditCaption edits captions of messages.
func (c Context) EditCaption(sig MessageSignature, cap CaptionData, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	cap.embed(&d)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.sig(sig, "editMessageCaption", d)
}

// EditMedia edits animation, audio, document, photo, or video messages.
func (c Context) EditMedia(sig MessageSignature, m tg.MediaInputter, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	err := prepareInputMedia(&d, false, m)
	if err != nil {
		return nil, err
	}
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.sig(sig, "editMessageMedia", d)
}

// EditReplyMarkup edits only the reply markup of messages.
func (c Context) EditReplyMarkup(sig MessageSignature, replyMarkup *tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData().SetJSON("reply_markup", replyMarkup)
	return c.sig(sig, "editMessageReplyMarkup", d)
}
