package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// MessageID creates a chat message id.
func MessageID(msg *tg.Message) MsgID {
	return MsgID{chatID: ChatID(msg.Chat.ID), msgID: msg.ID}
}

// InlineMsgID creates an inline message id.
func InlineMsgID(i *tg.InlineChosen) MsgID {
	return MsgID{inline: i.InlineMessageID}
}

// CallbackMsgID creates a MessageID of any callback message type.
func CallbackMsgID(q *tg.CallbackQuery) MsgID {
	if q.Message == nil {
		return MsgID{inline: q.InlineMessageID}
	}
	return MsgID{
		chatID: ChatID(q.Message.Chat().ID),
		msgID:  q.Message.ID(),
	}
}

// MsgID contains inline message id or chat id with message id.
type MsgID struct {
	chatID Chat
	msgID  int
	inline string
}

func (s MsgID) isInline() bool {
	return s.inline != ""
}

func (s MsgID) signature(d *api.Data) {
	if s.isInline() {
		d.Set("inline_message_id", s.inline)
		return
	}
	s.chatID.setChatID(d)
	d.SetInt("message_id", s.msgID)
}

// OpenMessage makes message interface.
func (c Context) OpenMessage(sig MsgID) MessageContext {
	return MessageContext{c, sig}
}

// MessageContext provides api for any message.
//
// For inline messages the result message will always be nil.
type MessageContext struct {
	Context
	sig MsgID
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
	LivePeriod         int64
	HorizontalAccuracy *float32
	Heading            int
	AlertRadius        int
}

// EditLiveLocation edits live location messages.
func (c MessageContext) EditLiveLocation(l LiveLocation, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	d.SetFloat("latitude", l.Lat)
	d.SetFloat("longitude", l.Long)
	d.SetInt64("live_period", l.LivePeriod)
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
	Text               string
	ParseMode          tg.ParseMode
	Entities           []tg.MessageEntity
	LinkPreviewOptions tg.LinkPreviewOptions
}

// EditText edits text and game messages.
func (c MessageContext) EditText(t EditText, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	d.Set("text", t.Text)
	d.Set("parse_mode", string(t.ParseMode))
	d.SetJSON("entities", t.Entities)
	d.SetJSON("link_preview_options", t.LinkPreviewOptions)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("editMessageText", d)
}

// EditCaption edits captions of messages.
func (c MessageContext) EditCaption(cap CaptionData, showCaptionAboveMedia bool, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	cap.embed(d)
	d.SetBool("show_caption_above_media", showCaptionAboveMedia)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("editMessageCaption", d)
}

// EditMedia edits animation, audio, document, photo, or video messages.
func (c MessageContext) EditMedia(m tg.MediaInputter, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	prepareInputMedia(d, m)
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

// SetReaction changes the chosen reactions on a message.
func (c MessageContext) SetReaction(reaction []tg.ReactionType, isBig ...bool) error {
	d := api.NewData()
	d.SetJSON("reaction", reaction)
	if len(isBig) > 0 {
		d.SetBool("is_big", isBig[0])
	}
	return c.method("setMessageReaction", d)
}
