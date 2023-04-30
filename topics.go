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

// EditForumTopic edits name and icon of a topic in a forum supergroup chat.
func (c ChatContext) EditForumTopic(threadID int, name, iconEmojiID string) error {
	d := api.NewData()
	d.Set("name", name)
	d.Set("icon_custom_emoji_id", iconEmojiID)
	return c.OpenForumTopic(threadID).method("editForumTopic", d)
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

// EditGeneralForumTopic edits the name of the 'General' topic in a forum supergroup chat.
func (c ChatContext) EditGeneralForumTopic(name string) error {
	d := api.NewData().Set("name", name)
	return c.method("editGeneralForumTopic", d)
}

// CloseGeneralForumTopic closes an open 'General' topic in a forum supergroup chat.
func (c ChatContext) CloseGeneralForumTopic() error {
	return c.method("closeGeneralForumTopic")
}

// ReopenGeneralForumTopic reopens a closed 'General' topic in a forum supergroup chat.
func (c ChatContext) ReopenGeneralForumTopic() error {
	return c.method("reopenGeneralForumTopic")
}

// HideGeneralForumTopic hides the 'General' topic in a forum supergroup chat.
func (c ChatContext) HideGeneralForumTopic() error {
	return c.method("hideGeneralForumTopic")
}

// UnhideGeneralForumTopic unhides the 'General' topic in a forum supergroup chat.
func (c ChatContext) UnhideGeneralForumTopic() error {
	return c.method("unhideGeneralForumTopic")
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
	d := api.NewData()
	mg.data(d)
	embed(d, opts)
	return topicMethod[[]tg.Message](t, "sendMediaGroup", d)
}

// SendChatAction sends chat action to tell the user that something
// is happening on the bot's side.
func (t Topic) SendChatAction(act tg.ChatAction) error {
	d := api.NewData().Set("action", string(act))
	return t.method("sendChatAction", d)
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
