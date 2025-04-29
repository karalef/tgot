package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// NewChatID makes ChatID from chat id.
func NewChatID(id int64, businness ...string) ChatID {
	return ChatID{id: id, business: firstOne(businness)}
}

// NewThreadID makes ChatID from thread id.
func NewThreadID(chatID, threadID int64, businness ...string) ChatID {
	return ChatID{thread: threadID, id: chatID, business: firstOne(businness)}
}

// Username makes ChatID from channel username.
func Username(username string, businness ...string) ChatID {
	return ChatID{username: username, business: firstOne(businness)}
}

func firstOne[T comparable](a []T, or ...T) (val T) {
	if len(a) > 0 {
		val = a[0]
	} else if len(or) > 0 {
		val = or[0]
	}
	return
}

// ChatID represents chat id or channel username.
type ChatID struct {
	id       int64
	thread   int64
	username string
	business string
}

func (c ChatID) ID() int64              { return c.id }
func (c ChatID) Username() string       { return c.username }
func (c ChatID) Thread() int64          { return c.thread }
func (c ChatID) BusinessConnID() string { return c.business }
func (c ChatID) IsTopic() bool          { return c.thread != 0 }
func (c ChatID) HasBusiness() bool      { return c.business != "" }

func (c ChatID) setChatID(d *api.Data, key ...string) {
	k := firstOne(key, "chat_id")
	if c.id != 0 {
		d.SetInt64(k, c.id)
	} else {
		d.Set(k, c.username)
	}
	if len(key) > 0 {
		return
	}
	d.SetInt64("message_thread_id", c.thread)
	d.Set("business_connection_id", c.BusinessConnID())
}

// WithChatID creates a new Chat from ctx and chat ID.
func WithChatID(ctx BaseContext, chatID ChatID) *Chat {
	d := api.NewData()
	chatID.setChatID(d)
	return &Chat{context: ctx.ctx().with(d), id: chatID}
}

type Chat struct {
	*context
	id ChatID
}

// WithName creates a new Chat with nested context name.
func (c *Chat) WithName(name string) *Chat {
	return &Chat{
		context: c.context.child(name),
		id:      c.id,
	}
}

// WithMember creates a new ChatMember with the specified user id.
func (c *Chat) WithMember(userID int64) *ChatMember {
	return WithChatMember(c, c.id, userID)
}

func (c *Chat) method(meth string, d ...*api.Data) error {
	_, err := method[bool](c.context, meth, d...)
	return err
}

// SendE sends the Sendable and returns only an error.
func (c *Chat) SendE(s Sendable, opts ...SendOptions) error {
	_, err := c.Send(s)
	return err
}

// ReplyTo creates ReplyParameters with only message ID.
func ReplyTo(msgID int) tg.ReplyParameters {
	return tg.ReplyParametersData[int64]{
		MessageID: msgID,
	}
}

// Reply replies to the message.
func (c *Chat) Reply(r tg.ReplyParameters, s Sendable) (*tg.Message, error) {
	return c.Send(s, SendOptions{ReplyParameters: r})
}

// ReplyE replies to the message and returns only an error.
func (c *Chat) ReplyE(r tg.ReplyParameters, s Sendable) error {
	return c.SendE(s, SendOptions{ReplyParameters: r})
}

func (c *Chat) ReplyTextSE(to int, text string, pm ...tg.ParseMode) error {
	msg := NewText(text)
	msg.ParseMode = firstOne(pm)
	return c.ReplyE(ReplyTo(to), msg)
}

// SendText sends just a text and returns only an error.
func (c *Chat) SendText(text string, pm ...tg.ParseMode) error {
	msg := Text{Text: text}
	if len(pm) > 0 {
		msg.ParseMode = pm[0]
	}
	return c.SendE(msg)
}

// GetInfo returns up to date information about the chat.
func (c *Chat) GetInfo() (*tg.ChatFullInfo, error) {
	return method[*tg.ChatFullInfo](c.context, "getChat")
}

// GetAdmins returns a list of administrators in a chat.
func (c *Chat) GetAdmins() ([]tg.ChatMember, error) {
	return method[[]tg.ChatMember](c.context, "getChatAdministrators")
}

// MemberCount returns the number of members in a chat.
func (c *Chat) MemberCount() (int, error) {
	return method[int](c.context, "getChatMemberCount")
}

// Leave a group, supergroup or channel.
func (c *Chat) Leave() error {
	return c.method("leaveChat")
}

// ForwardMany contains parameters for forwarding multiple messages.
type ForwardMany struct {
	MessageIDs          []int `tg:"message_ids"`
	DisableNotification bool  `tg:"disable_notification"`
	ProtectContent      bool  `tg:"protect_content"`
}

// ForwardMessages forwards multiple messages of any kind.
func (c *Chat) ForwardMessages(from ChatID, fwd ForwardMany) ([]tg.MessageID, error) {
	d := api.NewDataFrom(fwd)
	from.setChatID(d, "from_chat_id")
	return method[[]tg.MessageID](c.context, "forwardMessages", d)
}

// CopyMany contains parameters for copying multiple messages.
type CopyMany struct {
	MessageIDs    []int `tg:"message_ids"`
	RemoveCaption bool  `tg:"remove_caption"`
}

// CopyMessages copies messages of any kind.
func (c *Chat) CopyMessages(from ChatID, cp CopyMany, opts ...SendOptions) ([]tg.MessageID, error) {
	d := api.NewDataFrom(cp)
	from.setChatID(d, "from_chat_id")
	if len(opts) > 0 {
		d.SetObject(opts[0])
	}
	return method[[]tg.MessageID](c.context, "copyMessages", d)
}

// Send sends any Sendable object.
func (c *Chat) Send(s Sendable, opts ...SendOptions) (*tg.Message, error) {
	if s == nil {
		return nil, nil
	}
	d := api.NewDataFrom(s)
	if len(opts) > 0 {
		d.SetObject(opts[0])
	}
	return method[*tg.Message](c.context, s.sendMethod(), d)
}

// MediaGroup contains information about the media group to be sent.
type MediaGroup struct {
	Media       []tg.InputMedia          `tg:"media"`
	ReplyMarkup *tg.InlineKeyboardMarkup `tg:"reply_markup"`
}

// SendMediaGroup sends a group of photos, videos, documents or audios as an album.
func (c *Chat) SendMediaGroup(mg MediaGroup, opts ...SendOptions) ([]tg.Message, error) {
	d := api.NewDataFrom(mg)
	if len(opts) > 0 {
		d.SetObject(opts[0])
	}
	return method[[]tg.Message](c.context, "sendMediaGroup", d)
}

// SendChatAction sends chat action to tell the user that something
// is happening on the bot's side.
func (c *Chat) SendChatAction(act tg.ChatAction) error {
	d := api.NewData().Set("action", string(act))
	return c.method("sendChatAction", d)
}

// DeleteMessages deletes multiple messages simultaneously.
func (c *Chat) DeleteMessages(msgIDs []int) error {
	d := api.NewData().SetJSON("message_ids", msgIDs)
	return c.method("deleteMessages", d)
}

// BanSenderChat bans a channel chat in a supergroup or a channel.
func (c *Chat) BanSenderChat(senderID int64) error {
	d := api.NewData().SetInt64("sender_chat_id", senderID)
	return c.method("banChatSenderChat", d)
}

// UnbanSenderChat unbans a previously banned channel chat in a supergroup or channel.
func (c *Chat) UnbanSenderChat(senderID int64) error {
	d := api.NewData().SetInt64("sender_chat_id", senderID)
	return c.method("unbanChatSenderChat", d)
}

// SetPermissions sets default chat permissions for all members.
func (c *Chat) SetPermissions(perms tg.ChatPermissions, independentPerms ...bool) error {
	d := api.NewData().SetJSON("permissions", perms)
	d.SetBool("use_independent_chat_permissions", len(independentPerms) > 0 && independentPerms[0])
	return c.method("setChatPermissions", d)
}

// ExportInviteLink generates a new primary invite link for a chat;
// any previously generated primary link is revoked.
func (c *Chat) ExportInviteLink() (string, error) {
	return method[string](c.context, "exportChatInviteLink")
}

// InviteLink contains parameters for manipulations with invite links.
type InviteLink struct {
	Name               string
	ExpireDate         int64
	MemberLimit        int
	CreatesJoinRequest bool
}

func (i InviteLink) data() *api.Data {
	d := api.NewData()
	d.Set("name", i.Name)
	d.SetInt64("expire_date", i.ExpireDate)
	d.SetInt("member_limit", i.MemberLimit)
	d.SetBool("creates_join_request", i.CreatesJoinRequest)
	return d
}

// CreateInviteLink creates an additional invite link for a chat.
func (c *Chat) CreateInviteLink(i InviteLink) (*tg.ChatInviteLink, error) {
	return method[*tg.ChatInviteLink](c.context, "createChatInviteLink", i.data())
}

// EditInviteLink edits a non-primary invite link created by the bot.
func (c *Chat) EditInviteLink(link string, i InviteLink) (*tg.ChatInviteLink, error) {
	d := i.data().Set("invite_link", link)
	return method[*tg.ChatInviteLink](c.context, "editChatInviteLink", d)
}

// RevokeInviteLink revokes an invite link created by the bot.
func (c *Chat) RevokeInviteLink(link string) (*tg.ChatInviteLink, error) {
	d := api.NewData().Set("invite_link", link)
	return method[*tg.ChatInviteLink](c.context, "revokeChatInviteLink", d)
}

// SubscriptionInviteLink contains parameters for creating subscription invite links.
type SubscriptionInviteLink struct {
	Name   string
	Period uint
	Price  uint
}

func (i SubscriptionInviteLink) data() *api.Data {
	d := api.NewData()
	d.Set("name", i.Name)
	d.SetInt("period", int(i.Period))
	d.SetInt("price", int(i.Price))
	return d
}

// CreateSubscriptionInviteLink creates a subscription invite link for a channel chat.
func (c *Chat) CreateSubscriptionInviteLink(i SubscriptionInviteLink) (*tg.ChatInviteLink, error) {
	return method[*tg.ChatInviteLink](c.context, "createChatSubscriptionInviteLink", i.data())
}

// EditSubscriptionInviteLink edits a subscription invite link created by the bot.
func (c *Chat) EditSubscriptionInviteLink(link, name string) (*tg.ChatInviteLink, error) {
	d := api.NewData().Set("invite_link", link).Set("name", name)
	return method[*tg.ChatInviteLink](c.context, "editChatSubscriptionInviteLink", d)
}

// SetPhoto sets a new profile photo for the chat.
func (c *Chat) SetPhoto(photo *tg.InputFile) error {
	d := api.NewData().AddFile("photo", photo)
	return c.method("setChatPhoto", d)
}

// DeletePhoto deletes a chat photo.
func (c *Chat) DeletePhoto() error {
	return c.method("deleteChatPhoto")
}

// SetTitle change the title of a chat.
func (c *Chat) SetTitle(title string) error {
	d := api.NewData().Set("title", title)
	return c.method("setChatTitle", d)
}

// SetDescription changes the description of a group, a supergroup or a channel.
func (c *Chat) SetDescription(description string) error {
	d := api.NewData().Set("description", description)
	return c.method("setChatDescription", d)
}

// UnpinAllMessages clears the list of pinned messages in a chat.
func (c *Chat) UnpinAllMessages() error {
	return c.method("unpinAllChatMessages")
}

// SetStickerSet sets a new group sticker set for a supergroup.
func (c *Chat) SetStickerSet(stickerSet string) error {
	d := api.NewData().Set("sticker_set_name", stickerSet)
	return c.method("setChatStickerSet", d)
}

// DeleteStickerSet deletes a group sticker set from a supergroup.
func (c *Chat) DeleteStickerSet() error {
	return c.method("deleteChatStickerSet")
}

// SetMenuButton changes the bot's menu button in a private chat.
func (c *Chat) SetMenuButton(menu tg.MenuButton) error {
	d := api.NewData().SetJSON("menu_button", menu)
	return c.method("setChatMenuButton", d)
}

// GetMenuButton returns the current value of the bot's menu button in a private chat.
func (c *Chat) GetMenuButton() (*tg.MenuButton, error) {
	return method[*tg.MenuButton](c.context, "getChatMenuButton")
}

// CreateForumTopic creates a topic in a forum supergroup chat.
func (c *Chat) CreateForumTopic(name string, iconColor int, iconEmojiID string) (*tg.ForumTopic, error) {
	d := api.NewData().Set("name", name)
	d.SetInt("icon_color", iconColor)
	d.Set("icon_custom_emoji_id", iconEmojiID)
	return method[*tg.ForumTopic](c.context, "createForumTopic", d)
}

// EditGeneralForumTopic edits the name of the 'General' topic in a forum supergroup chat.
func (c *Chat) EditGeneralForumTopic(name string) error {
	d := api.NewData().Set("name", name)
	return c.method("editGeneralForumTopic", d)
}

// CloseGeneralForumTopic closes an open 'General' topic in a forum supergroup chat.
func (c *Chat) CloseGeneralForumTopic() error {
	return c.method("closeGeneralForumTopic")
}

// ReopenGeneralForumTopic reopens a closed 'General' topic in a forum supergroup chat.
func (c *Chat) ReopenGeneralForumTopic() error {
	return c.method("reopenGeneralForumTopic")
}

// HideGeneralForumTopic hides the 'General' topic in a forum supergroup chat.
func (c *Chat) HideGeneralForumTopic() error {
	return c.method("hideGeneralForumTopic")
}

// UnhideGeneralForumTopic unhides the 'General' topic in a forum supergroup chat.
func (c *Chat) UnhideGeneralForumTopic() error {
	return c.method("unhideGeneralForumTopic")
}

// UnpinAllGeneralMessages clears the list of pinned messages in a General forum topic.
func (c *Chat) UnpinAllGeneralForumTopicMessages() error {
	return c.method("unpinAllGeneralForumTopicMessages")
}

// EditForumTopic edits name and icon of a topic in a forum supergroup chat.
func (c *Chat) EditForumTopic(name, iconEmojiID string) error {
	d := api.NewData().Set("name", name)
	d.Set("icon_custom_emoji_id", iconEmojiID)
	return c.method("editForumTopic", d)
}

// CloseForumTopic closes an open topic in a forum supergroup chat.
func (c *Chat) CloseForumTopic() error {
	return c.method("closeForumTopic")
}

// ReopenForumTopic reopens a closed topic in a forum supergroup chat.
func (c *Chat) ReopenForumTopic() error {
	return c.method("reopenForumTopic")
}

// DeleteForumTopic deletes a forum topic along with all its messages in a forum supergroup chat.
func (c *Chat) DeleteForumTopic() error {
	return c.method("deleteForumTopic")
}

// UnpinAllMessages clears the list of pinned messages in a forum topic.
func (c *Chat) UnpinAllForumTopicMessages() error {
	return c.method("unpinAllForumTopicMessages")
}

// Verify verifies a chat on behalf of the organization which is represented by
// the bot.
func (c *Chat) Verify(desc string) error {
	return c.method("verifyChat", api.NewData().Set("custom_description", desc))
}

// RemoveVerification removes verification from a chat that is currently
// verified on behalf of the organization represented by the bot.
func (c *Chat) RemoveVerification() error {
	return c.method("removeChatVerification")
}
