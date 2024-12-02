package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// ChatMsgID creates a chat message id.
func ChatMsgID(msg *tg.Message) MessageID {
	return MessageID{chatID: NewChatID(msg.Chat.ID, msg.BusinessConnectionID), msgID: msg.ID}
}

// InlineMsgID creates an inline message id.
func InlineMsgID(i *tg.InlineChosen) MessageID {
	return MessageID{inline: i.InlineMessageID}
}

// CallbackMsgID creates a message id of any callback message type.
func CallbackMsgID(q *tg.CallbackQuery) MessageID {
	if q.Message == nil || q.Message.IsInaccessible() {
		return MessageID{inline: q.InlineMessageID}
	}
	return MessageID{
		chatID: NewChatID(q.Message.Chat().ID, q.Message.BusinessConnectionID),
		msgID:  q.Message.ID(),
	}
}

// MessageID contains inline message id or chat id with message id.
type MessageID struct {
	chatID ChatID
	msgID  int
	inline string
}

func (s MessageID) ChatID() ChatID { return s.chatID }
func (s MessageID) MessageID() int { return s.msgID }
func (s MessageID) Inline() string { return s.inline }

func (s MessageID) isInline() bool { return s.inline != "" }

func (s MessageID) setTo(d *api.Data) *api.Data {
	if s.isInline() {
		d.Set("inline_message_id", s.inline)
		return d
	}
	s.chatID.setChatID(d)
	d.SetInt("message_id", s.msgID)
	return d
}

func WithMessage(ctx BaseContext, id MessageID) *Message {
	return &Message{context: ctx.ctx().with(id.setTo(api.NewData())), sig: id}
}

// Message provides api for any message.
//
// For inline messages the result message will always be nil.
type Message struct {
	*context
	sig MessageID
}

// Chat returns the Chat of the Message.
// Returns nil for inline messages.
func (m *Message) Chat() *Chat {
	if m.sig.isInline() {
		return nil
	}
	return WithChatID(m, m.sig.chatID)
}

func (m *Message) WithName(name string) *Message {
	return &Message{context: m.context.child(name), sig: m.sig}
}

func (m *Message) ID() MessageID  { return m.sig }
func (m *Message) ChatID() ChatID { return m.sig.chatID }
func (m *Message) MessageID() int { return m.sig.msgID }
func (m *Message) Inline() string { return m.sig.inline }

func (c Message) msgMethod(meth string, d *api.Data) (*tg.Message, error) {
	if c.sig.isInline() {
		return nil, c.method(meth, d)
	}
	return method[*tg.Message](c, meth, d)
}

// ReplyText replies to the message with text.
func (m *Message) ReplyText(text string, pm ...tg.ParseMode) error {
	return m.Chat().ReplyTextSE(m.MessageID(), text, pm...)
}

// Reply replies to the message.
func (m *Message) Reply(s Sendable) error {
	return m.Chat().ReplyE(ReplyTo(m.MessageID()), s)
}

// EditLiveLocation edits live location messages.
func (c Message) EditLiveLocation(ll tg.Location, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	api.MarshalTo(d, ll, "json")
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("editMessageLiveLocation", d)
}

// StopLiveLocation stops updating a live location message before live_period expires.
func (c Message) StopLiveLocation(replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("stopMessageLiveLocation", d)
}

// EditText contains parameters for editing message text.
type EditText struct {
	Text               string                `tg:"text"`
	ParseMode          tg.ParseMode          `tg:"parse_mode"`
	Entities           []tg.MessageEntity    `tg:"entities"`
	LinkPreviewOptions tg.LinkPreviewOptions `tg:"link_preview_options"`
}

// EditText edits text and game messages.
func (c Message) EditText(t EditText, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewDataFrom(t)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("editMessageText", d)
}

// EditCaption contains parameters for editing message caption.
type EditCaption struct {
	CaptionData
	ShowCaptionAboveMedia bool `tg:"show_caption_above_media"`
}

// EditCaption edits captions of messages.
func (c Message) EditCaption(cap EditCaption, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewDataFrom(cap)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("editMessageCaption", d)
}

// EditMedia edits animation, audio, document, photo, or video messages.
func (c Message) EditMedia(m tg.InputMedia, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData()
	med, thumb := m.GetMedia()
	if med != nil {
		d.AddAttach(med)
		d.AddAttach(thumb)
	}
	d.SetJSON("media", m)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return c.msgMethod("editMessageMedia", d)
}

// EditReplyMarkup edits only the reply markup of messages.
func (c Message) EditReplyMarkup(replyMarkup *tg.InlineKeyboardMarkup) (*tg.Message, error) {
	d := api.NewData().SetJSON("reply_markup", replyMarkup)
	return c.msgMethod("editMessageReplyMarkup", d)
}

// SetReaction changes the chosen reactions on a message.
func (c Message) SetReaction(reaction []tg.ReactionType, isBig ...bool) error {
	d := api.NewData().SetJSON("reaction", reaction)
	if len(isBig) > 0 {
		d.SetBool("is_big", isBig[0])
	}
	return c.method("setMessageReaction", d)
}

// StopPoll stops a poll which was sent by the bot.
func (c Message) StopPoll(replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Poll, error) {
	d := api.NewData()
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return method[*tg.Poll](c, "stopPoll", d)
}

// DeleteMessage deletes a message, including service messages.
func (c Message) Delete() error {
	return c.method("deleteMessage")
}

// Copy contains parameters for copying the message.
type Copy struct {
	CaptionData
	ReplyMarkup           tg.ReplyMarkup `tg:"reply_markup"`
	ShowCaptionAboveMedia bool           `tg:"show_caption_above_media"`
}

// Copy copies messages of any kind.
// Service messages and invoice messages can't be copied.
func (c Message) Copy(from ChatID, cp Copy, opts ...SendOptions) (*tg.MessageID, error) {
	d := api.NewDataFrom(cp)
	from.setChatID(d, "from_chat_id")
	if len(opts) > 0 {
		d.SetObject(opts[0])
	}
	return method[*tg.MessageID](c, "copyMessage", d)
}

// Forward contains parameters for forwarding the message.
type Forward struct {
	DisableNotification bool `tg:"disable_notification"`
	ProtectContent      bool `tg:"protect_content"`
}

// Forward forwards messages of any kind.
// Service messages can't be forwarded.
func (c Message) Forward(from ChatID, fwd Forward) (*tg.Message, error) {
	d := api.NewDataFrom(fwd)
	from.setChatID(d, "from_chat_id")
	return c.msgMethod("forwardMessage", d)
}

// SetGameScore contains parameters for setting the game score.
type SetGameScore struct {
	Score       int  `tg:"score,force"`
	Force       bool `tg:"force"`
	DisableEdit bool `tg:"disable_edit_message"`
}

// SetGameScore sets the score of the specified user in a game message.
func (c Message) SetGameScore(userID int64, s SetGameScore) (*tg.Message, error) {
	d := api.NewData()
	d.SetInt64("user_id", userID)
	d.SetInt("score", s.Score, true)
	d.SetBool("force", s.Force)
	d.SetBool("disable_edit_message", s.DisableEdit)
	return c.msgMethod("setGameScore", d)
}

// GetGameHighScores returns data for high score tables.
func (c Message) GetGameHighScores(userID int64) ([]tg.GameHighScore, error) {
	d := api.NewData().SetInt64("user_id", userID)
	return method[[]tg.GameHighScore](c, "getGameHighScores", d)
}

// PinMessage adds a message to the list of pinned messages in a chat.
func (c Message) PinMessage(notify bool) error {
	d := api.NewData().SetBool("disable_notification", !notify)
	return c.method("pinChatMessage", d)
}

// UnpinMessage removes a message from the list of pinned messages in a chat.
func (c Message) UnpinMessage() error {
	return c.method("unpinChatMessage")
}
