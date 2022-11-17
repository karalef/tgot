package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// GetForumTopicIconStickers returns custom emoji stickers,
// which can be used as a forum topic icon by any user.
func (c Context) GetForumTopicIconStickers() ([]tg.Sticker, error) {
	return method[[]tg.Sticker](c, "getForumTopicIconStickers")
}

// Topic provides forum topics api.
type Topic struct {
	chat     *Chat
	threadID int
}

func (t *Topic) method(method string, d ...api.Data) error {
	_, err := topicMethod[bool](t, method, d...)
	return err
}

func topicMethod[T any](t *Topic, meth string, d ...api.Data) (T, error) {
	var data api.Data
	if len(d) > 0 {
		data = d[0]
	} else {
		data = api.NewData()
	}
	t.chat.setChatID(&data)
	data.SetInt("message_thread_id", t.threadID)
	return method[T](t.chat.Context, meth, data)
}

// Edit edits name and icon of a topic in a forum supergroup chat.
func (t *Topic) Edit(name, iconEmojiID string) error {
	d := api.NewData()
	d.Set("name", name)
	d.Set("icon_custom_emoji_id", iconEmojiID)
	return t.method("editForumTopic", d)
}

// Close closes an open topic in a forum supergroup chat.
func (t *Topic) Close() error {
	return t.method("closeForumTopic")
}

// Reopen reopens a closed topic in a forum supergroup chat.
func (t *Topic) Reopen() error {
	return t.method("reopenForumTopic")
}

// Delete deletes a forum topic along with all its messages in a forum supergroup chat.
func (t *Topic) Delete() error {
	return t.method("deleteForumTopic")
}

// UnpinAllMessages clears the list of pinned messages in a forum topic.
func (t *Topic) UnpinAllMessages() error {
	return t.method("unpinAllForumTopicMessages")
}

// Send sends any Sendable object.
func (t *Topic) Send(s Sendable, opts ...SendOptions[tg.ReplyMarkup]) (*tg.Message, error) {
	if s == nil {
		return nil, nil
	}

	method, d := s.data()
	if len(opts) > 0 {
		opts[0].embed(&d)
	}
	return topicMethod[*tg.Message](t, method, d)
}

// SendMediaGroup sends a group of photos, videos, documents or audios as an album.
func (t *Topic) SendMediaGroup(mg MediaGroup, opts ...BaseSendOptions) ([]tg.Message, error) {
	d, err := mg.data()
	if err != nil {
		return nil, err
	}
	if len(opts) > 0 {
		opts[0].embed(&d)
	}
	return topicMethod[[]tg.Message](t, "sendMediaGroup", d)
}

// SendInvoice sends an invoice.
func (t *Topic) SendInvoice(i Invoice, opts ...SendOptions[*tg.InlineKeyboardMarkup]) (*tg.Message, error) {
	d := i.data()
	if len(opts) > 0 {
		opts[0].embed(&d)
	}
	return topicMethod[*tg.Message](t, "sendInvoice", d)
}

// SendGame sends a game.
func (t *Topic) SendGame(g Game, opts ...SendOptions[*tg.InlineKeyboardMarkup]) (*tg.Message, error) {
	d := g.data()
	if len(opts) > 0 {
		opts[0].embed(&d)
	}
	return topicMethod[*tg.Message](t, "sendGame", d)
}

// Forward forwards messages of any kind.
// Service messages can't be forwarded.
func (t *Topic) Forward(from *Chat, fwd Forward) (*tg.Message, error) {
	d := api.NewData()
	from.setChatID(&d, "from_chat_id")
	d.SetInt("message_id", fwd.MessageID)
	d.SetBool("disable_notification", fwd.DisableNotification)
	d.SetBool("protect_content", fwd.ProtectContent)
	return topicMethod[*tg.Message](t, "forwardMessage", d)
}

// Copy copies messages of any kind.
// Service messages and invoice messages can't be copied.
func (t *Topic) Copy(from *Chat, cp Copy, opts ...SendOptions[tg.ReplyMarkup]) (*tg.MessageID, error) {
	d := api.NewData()
	from.setChatID(&d, "from_chat_id")
	d.SetInt("message_id", cp.MessageID)
	cp.CaptionData.embed(&d)
	if len(opts) > 0 {
		opts[0].embed(&d)
	}
	return topicMethod[*tg.MessageID](t, "copyMessage", d)
}
