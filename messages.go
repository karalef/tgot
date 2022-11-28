package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// MessageSignature creates a chat message signature.
func MessageSignature(msg *tg.Message) MsgSignature {
	return MsgSignature{chatID: ChatID(msg.Chat.ID)}
}

// InlineSignature creates an inline message signature.
func InlineSignature(i *tg.InlineChosen) MsgSignature {
	return MsgSignature{inline: i.InlineMessageID}
}

// CallbackSignature creates a MessageSignature of any callback message type.
func CallbackSignature(q *tg.CallbackQuery) MsgSignature {
	if q.Message != nil {
		return MessageSignature(q.Message)
	}
	return MsgSignature{inline: q.InlineMessageID}
}

// MsgSignature contains inline message id or chat id with message id.
type MsgSignature struct {
	chatID Chat
	msgID  int
	inline string
}

func (s MsgSignature) isInline() bool {
	return s.inline != ""
}

func (s MsgSignature) signature(d *api.Data) {
	if s.isInline() {
		d.Set("inline_message_id", s.inline)
		return
	}
	s.chatID.setChatID(d)
	d.SetInt("message_id", s.msgID)
}

// OpenMessage makes message interface.
func (c Context) OpenMessage(sig MsgSignature) MessageContext {
	return MessageContext{c, sig}
}

// MessageContext provides api for any message.
//
// For inline messages the result message will always be nil.
type MessageContext struct {
	Context
	sig MsgSignature
}

func (c MessageContext) msgMethod(meth string, d *api.Data) (*tg.Message, error) {
	c.sig.signature(d)
	if c.sig.isInline() {
		return nil, c.method(meth, d)
	}
	return method[*tg.Message](c.Context, meth, d)
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
func (c MessageContext) EditLiveLocation(l LiveLocation, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
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
	return c.msgMethod("editMessageLiveLocation", d)
}

// StopLiveLocation stops updating a live location message before live_period expires.
func (c MessageContext) StopLiveLocation(replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("stopMessageLiveLocation", d)
}

// EditText contains parameters for editing message text.
type EditText struct {
	Text                  string
	ParseMode             tg.ParseMode
	Entities              []tg.MessageEntity
	DisableWebPagePreview bool
}

// EditText edits text and game messages.
func (c MessageContext) EditText(t EditText, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	d.Set("text", t.Text)
	d.Set("parse_mode", string(t.ParseMode))
	d.SetJSON("entities", t.Entities)
	d.SetBool("disable_web_page_preview", t.DisableWebPagePreview)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("editMessageText", d)
}

// EditCaption edits captions of messages.
func (c MessageContext) EditCaption(cap CaptionData, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	cap.embed(d)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("editMessageCaption", d)
}

// EditMedia edits animation, audio, document, photo, or video messages.
func (c MessageContext) EditMedia(m tg.MediaInputter, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	err := prepareInputMedia(d, m)
	if err != nil {
		return nil, err
	}
	d.SetJSON("media", m)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("editMessageMedia", d)
}

// EditReplyMarkup edits only the reply markup of messages.
func (c MessageContext) EditReplyMarkup(replyMarkup *tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData().SetJSON("reply_markup", replyMarkup)
	return c.msgMethod("editMessageReplyMarkup", d)
}
