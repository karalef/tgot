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

// OpenForumTopic makes forum topic interface.
func (c ChatContext) OpenForumTopic(threadID int) Topic {
	return Topic{c.Context, c.id, threadID}
}

// CreateForumTopic creates a topic in a forum supergroup chat.
func (c ChatContext) CreateForumTopic(name string, iconColor int, iconEmojiID string) (*tg.ForumTopic, error) {
	d := api.NewData().Set("name", name)
	d.SetInt("icon_color", iconColor)
	d.Set("icon_custom_emoji_id", iconEmojiID)
	return chatMethod[*tg.ForumTopic](c, "createForumTopic", d)
}

// CloseForumTopic closes an open topic in a forum supergroup chat.
func (c ChatContext) CloseForumTopic(threadID int) error {
	return c.OpenForumTopic(threadID).method("closeForumTopic")
}

// ReopenForumTopic reopens a closed topic in a forum supergroup chat.
func (c ChatContext) ReopenForumTopic(threadID int) error {
	return c.OpenForumTopic(threadID).method("reopenForumTopic")
}

// DeleteForumTopic deletes a forum topic along with all its messages in a forum supergroup chat.
func (c ChatContext) DeleteForumTopic(threadID int) error {
	return c.OpenForumTopic(threadID).method("deleteForumTopic")
}

// Topic provides forum topics api.
type Topic struct {
	ctx      Context
	chatID   Chat
	threadID int
}

func (t Topic) method(method string, d ...*api.Data) error {
	_, err := topicMethod[bool](t, method, d...)
	return err
}

func topicMethod[T any](t Topic, meth string, d ...*api.Data) (T, error) {
	var data *api.Data
	if len(d) > 0 {
		data = d[0]
	} else {
		data = api.NewData()
	}
	data.SetInt("message_thread_id", t.threadID)
	t.chatID.setChatID(data)
	return method[T](t.ctx, meth, data)
}

// Edit edits name and icon of a topic in a forum supergroup chat.
func (t Topic) Edit(name, iconEmojiID string) error {
	d := api.NewData()
	d.Set("name", name)
	d.Set("icon_custom_emoji_id", iconEmojiID)
	return t.method("editForumTopic", d)
}

// UnpinAllMessages clears the list of pinned messages in a forum topic.
func (t Topic) UnpinAllMessages() error {
	return t.method("unpinAllForumTopicMessages")
}

// Send sends any Sendable object.
func (t Topic) Send(s Sendable, opts ...SendOptions) (*tg.Message, error) {
	if s == nil {
		return nil, nil
	}

	d := api.NewData()
	embed(d, opts)
	return topicMethod[*tg.Message](t, s.sendData(d), d)
}

// SendMediaGroup sends a group of photos, videos, documents or audios as an album.
func (t Topic) SendMediaGroup(mg MediaGroup, opts ...SendOptions) ([]tg.Message, error) {
	d, err := mg.data()
	if err != nil {
		return nil, err
	}
	embed(d, opts)
	return topicMethod[[]tg.Message](t, "sendMediaGroup", d)
}

// Forward forwards messages of any kind.
// Service messages can't be forwarded.
func (t Topic) Forward(from Chat, fwd Forward) (*tg.Message, error) {
	d := api.NewData()
	from.setChatID(d, "from_chat_id")
	fwd.data(d)
	return topicMethod[*tg.Message](t, "forwardMessage", d)
}

// Copy copies messages of any kind.
// Service messages and invoice messages can't be copied.
func (t Topic) Copy(from Chat, cp Copy, opts ...SendOptions) (*tg.MessageID, error) {
	d := api.NewData()
	from.setChatID(d, "from_chat_id")
	cp.data(d)
	if len(opts) > 0 {
		opts[0].embed(d)
	}
	return topicMethod[*tg.MessageID](t, "copyMessage", d)
}
