package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/tg"
)

func (b *Bot) makeChatContext(chat *tg.Chat, name string) ChatContext {
	return ChatContext{Chat: b.MakeContext(name).OpenChat(chat.ID)}
}

// ChatContext type.
type ChatContext struct {
	Chat
}

// Send sends the Sendable and returns only an error.
// It is short version for [ChatContext.Chat.Send].
func (c ChatContext) Send(s Sendable, opts ...SendOptions[tg.ReplyMarkup]) error {
	_, err := c.Chat.Send(s)
	return err
}

// SendText sends just a text and returns only an error.
func (c ChatContext) SendText(text string, pm ...tg.ParseMode) error {
	msg := Message{Text: text}
	if len(pm) > 0 {
		msg.ParseMode = pm[0]
	}
	return c.Send(msg)
}

// Chat represents chat api.
type Chat struct {
	Context
	chatID   int64
	username string
}

func (c *Chat) setChatID(d *api.Data, key ...string) {
	k := "chat_id"
	if len(key) > 0 {
		k = key[0]
	}
	if c.chatID != 0 {
		d.SetInt64(k, c.chatID)
	} else {
		d.Set(k, c.username)
	}
}

func (c *Chat) method(method string, d ...api.Data) error {
	_, err := chatMethod[bool](c, method, d...)
	return err
}

func chatMethod[T any](c *Chat, meth string, d ...api.Data) (T, error) {
	var data api.Data
	if len(d) > 0 {
		data = d[0]
	} else {
		data = api.NewData()
	}

	c.setChatID(&data)
	return method[T](c.Context, meth, data)
}

// GetInfo returns up to date information about the chat.
func (c *Chat) GetInfo() (*tg.Chat, error) {
	return chatMethod[*tg.Chat](c, "getChat")
}

// GetAdmins returns a list of administrators in a chat.
func (c *Chat) GetAdmins() ([]tg.ChatMember, error) {
	return chatMethod[[]tg.ChatMember](c, "getChatAdministrators")
}

// MemberCount returns the number of members in a chat.
func (c *Chat) MemberCount() (int, error) {
	return chatMethod[int](c, "getChatMemberCount")
}

// GetMember returns information about a member of a chat.
func (c *Chat) GetMember(userID int64) (*tg.ChatMember, error) {
	d := api.NewData().SetInt64("user_id", userID)
	return chatMethod[*tg.ChatMember](c, "getChatMember", d)
}

// Leave a group, supergroup or channel.
func (c *Chat) Leave() error {
	return c.method("leaveChat")
}

// Forward contains paramenets for forwarding the message.
type Forward struct {
	MessageID           int
	DisableNotification bool
	ProtectContent      bool
}

// Forward forwards messages of any kind.
// Service messages can't be forwarded.
func (c *Chat) Forward(from *Chat, fwd Forward) (*tg.Message, error) {
	d := api.NewData()
	from.setChatID(&d, "from_chat_id")
	d.SetInt("message_id", fwd.MessageID)
	d.SetBool("disable_notification", fwd.DisableNotification)
	d.SetBool("protect_content", fwd.ProtectContent)
	return chatMethod[*tg.Message](c, "forwardMessage", d)
}

// ForwardTo forwards the message to the specified chat instead of current.
func (c *Chat) ForwardTo(to *Chat, fwd Forward) (*tg.Message, error) {
	return to.Forward(c, fwd)
}

// Copy contains parameters for copying the message.
type Copy struct {
	MessageID int
	CaptionData
}

// Copy copies messages of any kind.
// Service messages and invoice messages can't be copied.
func (c *Chat) Copy(from *Chat, cp Copy, opts ...SendOptions[tg.ReplyMarkup]) (*tg.MessageID, error) {
	d := api.NewData()
	from.setChatID(&d, "from_chat_id")
	d.SetInt("message_id", cp.MessageID)
	cp.CaptionData.embed(&d)
	if len(opts) > 0 {
		opts[0].embed(&d)
	}
	return chatMethod[*tg.MessageID](c, "copyMessage", d)
}

// CopyTo copies the message to the specified chat instead of current.
func (c *Chat) CopyTo(to *Chat, cp Copy, opts ...SendOptions[tg.ReplyMarkup]) (*tg.MessageID, error) {
	return to.Copy(c, cp)
}

// SendText sends just a text.
func (c *Chat) SendText(text string) (*tg.Message, error) {
	return c.Send(NewMessage(text))
}

// Send sends any Sendable object.
func (c *Chat) Send(s Sendable, opts ...SendOptions[tg.ReplyMarkup]) (*tg.Message, error) {
	if s == nil {
		return nil, nil
	}

	method, d := s.data()
	if len(opts) > 0 {
		opts[0].embed(&d)
	}

	return chatMethod[*tg.Message](c, method, d)
}

// SendMediaGroup sends a group of photos, videos, documents or audios as an album.
func (c *Chat) SendMediaGroup(mg MediaGroup, opts ...BaseSendOptions) ([]tg.Message, error) {
	d, err := mg.data()
	if err != nil {
		return nil, err
	}
	if len(opts) > 0 {
		opts[0].embed(&d)
	}
	return chatMethod[[]tg.Message](c, "sendMediaGroup", d)
}

// SendChatAction sends chat action to tell the user that something
// is happening on the bot's side.
func (c *Chat) SendChatAction(act tg.ChatAction) error {
	d := api.NewData().Set("action", string(act))
	return c.method("sendChatAction", d)
}

// SendInvoice sends an invoice.
func (c *Chat) SendInvoice(i Invoice, opts ...SendOptions[*tg.InlineKeyboardMarkup]) (*tg.Message, error) {
	d := i.data()
	if len(opts) > 0 {
		opts[0].embed(&d)
	}
	return chatMethod[*tg.Message](c, "sendInvoice", d)
}

// SendGame sends a game.
func (c *Chat) SendGame(g Game, opts ...SendOptions[*tg.InlineKeyboardMarkup]) (*tg.Message, error) {
	d := g.data()
	if len(opts) > 0 {
		opts[0].embed(&d)
	}
	return chatMethod[*tg.Message](c, "sendGame", d)
}

// StopPoll stops a poll which was sent by the bot.
func (c *Chat) StopPoll(msgID int, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Poll, error) {
	d := api.NewData()
	d.SetInt("message_id", msgID)
	if len(replyMarkup) > 0 {
		d.SetJSON("reply_markup", replyMarkup[0])
	}
	return chatMethod[*tg.Poll](c, "stopPoll", d)
}

// DeleteMessage deletes a message, including service messages.
func (c *Chat) DeleteMessage(msgID int) error {
	d := api.NewData().SetInt("message_id", msgID)
	return c.method("deleteMessage", d)
}

// Ban contains parameters for banning a chat member.
type Ban struct {
	UserID         int64
	UntilDate      *int64
	RevokeMessages bool
}

// Ban bans a user in a group, a supergroup or a channel.
func (c *Chat) Ban(b Ban) error {
	d := api.NewData()
	d.SetInt64("user_id", b.UserID)
	if b.UntilDate != nil {
		d.SetInt64("until_date", *b.UntilDate, true)
	}
	d.SetBool("revoke_messages", b.RevokeMessages)
	return c.method("banChatMember", d)
}

// Unban unbans a previously banned user in a supergroup or channel.
func (c *Chat) Unban(userID int64, onlyIfBanned bool) error {
	d := api.NewData().SetInt64("user_id", userID)
	d.SetBool("only_if_banned", onlyIfBanned)
	return c.method("unbanChatMember", d)
}

// Restrict restricts a user in a supergroup.
func (c *Chat) Restrict(userID int64, perms tg.ChatPermissions, until *int64) error {
	d := api.NewData().SetInt64("user_id", userID)
	d.SetJSON("permissions", perms)
	if until != nil {
		d.SetInt64("until_date", *until, true)
	}
	return c.method("restrictChatMember", d)
}

// Promote promotes or demotes a user in a supergroup or a channel.
func (c *Chat) Promote(userID int64, rights tg.ChatAdministratorRights) error {
	d := api.NewData().SetInt64("user_id", userID)
	d.SetBool("is_anonymous", rights.IsAnonymous)
	d.SetBool("can_manage_chat", rights.CanManageChat)
	d.SetBool("can_delete_messages", rights.CanDeleteMessages)
	d.SetBool("can_manage_video_chats", rights.CanManageVideoChats)
	d.SetBool("can_restrict_members", rights.CanRestrictMembers)
	d.SetBool("can_promote_members", rights.CanPromoteMembers)
	d.SetBool("can_change_info", rights.CanChangeInfo)
	d.SetBool("can_invite_users", rights.CanInviteUsers)
	d.SetBool("can_post_messages", rights.CanPostMessages)
	d.SetBool("can_edit_messages", rights.CanEditMessages)
	d.SetBool("can_pin_messages", rights.CanPinMessages)
	return c.method("promoteChatMember", d)
}

// SetAdminTitle sets a custom title for an administrator in a supergroup promoted by the bot.
func (c *Chat) SetAdminTitle(userID int64, title string) error {
	d := api.NewData()
	d.SetInt64("user_id", userID)
	d.Set("custom_title", title)
	return c.method("setChatAdministratorCustomTitle", d)
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
func (c *Chat) SetPermissions(perms tg.ChatPermissions) error {
	d := api.NewData().SetJSON("permissions", perms)
	return c.method("setChatPermissions", d)
}

// ExportInviteLink generates a new primary invite link for a chat;
// any previously generated primary link is revoked.
func (c *Chat) ExportInviteLink() (string, error) {
	return chatMethod[string](c, "exportChatInviteLink")
}

// InviteLink contains parameters for manipulations with invite links.
type InviteLink struct {
	Name               string
	ExpireDate         int64
	MemberLimit        int
	CreatesJoinRequest bool
}

func (i InviteLink) data() api.Data {
	d := api.NewData()
	d.Set("name", i.Name)
	d.SetInt64("expire_date", i.ExpireDate)
	d.SetInt("member_limit", i.MemberLimit)
	d.SetBool("creates_join_request", i.CreatesJoinRequest)
	return d
}

// CreateInviteLink creates an additional invite link for a chat.
func (c *Chat) CreateInviteLink(i InviteLink) (*tg.ChatInviteLink, error) {
	return chatMethod[*tg.ChatInviteLink](c, "createChatInviteLink", i.data())
}

// EditInviteLink edits a non-primary invite link created by the bot.
func (c *Chat) EditInviteLink(link string, i InviteLink) (*tg.ChatInviteLink, error) {
	d := i.data().Set("invite_link", link)
	return chatMethod[*tg.ChatInviteLink](c, "editChatInviteLink", d)
}

// RevokeInviteLink revokes an invite link created by the bot.
func (c *Chat) RevokeInviteLink(link string) (*tg.ChatInviteLink, error) {
	d := api.NewData().Set("invite_link", link)
	return chatMethod[*tg.ChatInviteLink](c, "revokeChatInviteLink", d)
}

// ApproveJoinRequest approves a chat join request.
func (c *Chat) ApproveJoinRequest(userID int64) error {
	d := api.NewData().SetInt64("user_id", userID)
	return c.method("approveChatJoinRequest", d)
}

// DeclineJoinRequest declines a chat join request.
func (c *Chat) DeclineJoinRequest(userID int64) error {
	d := api.NewData().SetInt64("user_id", userID)
	return c.method("declineChatJoinRequest", d)
}

// SetPhoto sets a new profile photo for the chat.
func (c *Chat) SetPhoto(photo *tg.InputFile) error {
	d := api.NewData()
	d.SetFile("photo", photo, nil)
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

// PinMessage adds a message to the list of pinned messages in a chat.
func (c *Chat) PinMessage(msgID int, notify bool) error {
	d := api.NewData().SetInt("message_id", msgID)
	d.SetBool("disable_notification", !notify)
	return c.method("pinChatMessage", d)
}

// UnpinMessage removes a message from the list of pinned messages in a chat.
func (c *Chat) UnpinMessage(msgID int) error {
	d := api.NewData().SetInt("message_id", msgID)
	return c.method("unpinChatMessage", d)
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
	return chatMethod[*tg.MenuButton](c, "getChatMenuButton")
}
