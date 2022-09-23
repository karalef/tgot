package bot

import (
	"tghwbot/bot/tg"
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
func (c ChatContext) Send(s Sendable) error {
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

func (c *Chat) setChatID(p params, key ...string) {
	k := "chat_id"
	if len(key) > 0 {
		k = key[0]
	}
	if c.chatID != 0 {
		p.setInt64(k, c.chatID)
	} else {
		p.set(k, c.username)
	}
}

func (c *Chat) method(method string, p params, files ...file) error {
	_, err := chatMethod[bool](c, method, p, files...)
	return err
}

func chatMethod[T any](c *Chat, method string, p params, files ...file) (T, error) {
	if p == nil {
		p = params{}
	}

	c.setChatID(p)
	return api[T](c.Context, method, p, files...)
}

// GetInfo returns up to date information about the chat.
func (c *Chat) GetInfo() (*tg.Chat, error) {
	return chatMethod[*tg.Chat](c, "getChat", nil)
}

// GetAdmins returns a list of administrators in a chat.
func (c *Chat) GetAdmins() ([]tg.ChatMember, error) {
	return chatMethod[[]tg.ChatMember](c, "getChatAdministrators", nil)
}

// MemberCount returns the number of members in a chat.
func (c *Chat) MemberCount() (int, error) {
	return chatMethod[int](c, "getChatMemberCount", nil)
}

// GetMember returns information about a member of a chat.
func (c *Chat) GetMember(userID int64) (*tg.ChatMember, error) {
	p := params{}.setInt64("user_id", userID)
	return chatMethod[*tg.ChatMember](c, "getChatMember", p)
}

// Leave a group, supergroup or channel.
func (c *Chat) Leave() error {
	return c.method("leaveChat", nil)
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
	p := params{}
	from.setChatID(p, "from_chat_id")
	p.setInt("message_id", fwd.MessageID)
	p.setBool("disable_notification", fwd.DisableNotification)
	p.setBool("protect_content", fwd.ProtectContent)
	return chatMethod[*tg.Message](c, "forwardMessage", p)
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
func (c *Chat) Copy(from *Chat, cp Copy, opts ...SendOptions[tg.ReplyMarkup]) (*tg.Message, error) {
	p := params{}
	from.setChatID(p, "from_chat_id")
	p.setInt("message_id", cp.MessageID)
	cp.CaptionData.embed(p)
	if len(opts) > 0 {
		opts[0].embed(p)
	}
	return chatMethod[*tg.Message](c, "copyMessage", p)
}

// CopyTo copies the message to the specified chat instead of current.
func (c *Chat) CopyTo(to *Chat, cp Copy, opts ...SendOptions[tg.ReplyMarkup]) (*tg.Message, error) {
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

	p := params{}
	files := s.data(p)
	if len(opts) > 0 {
		opts[0].embed(p)
	}

	return chatMethod[*tg.Message](c, s.method(), p, files...)
}

// SendMediaGroup sends a group of photos, videos, documents or audios as an album.
func (c *Chat) SendMediaGroup(mg MediaGroup, opts ...MediaGroupSendOptions) ([]tg.Message, error) {
	p := params{}
	files, err := mg.data(p)
	if err != nil {
		return nil, err
	}
	if len(opts) > 0 {
		opts[0].embed(p)
	}
	return chatMethod[[]tg.Message](c, "sendMediaGroup", p, files...)
}

// SendChatAction sends chat action to tell the user that something
// is happening on the bot's side.
func (c *Chat) SendChatAction(act tg.ChatAction) error {
	p := params{}.set("action", string(act))
	return c.method("sendChatAction", p)
}

// SendInvoice sends an invoice.
func (c *Chat) SendInvoice(i Invoice, opts ...SendOptions[*tg.InlineKeyboardMarkup]) (*tg.Message, error) {
	p := params{}
	i.params(p)
	if len(opts) > 0 {
		opts[0].embed(p)
	}
	return chatMethod[*tg.Message](c, "sendInvoice", p)
}

// SendGame sends a game.
func (c *Chat) SendGame(g Game, opts ...SendOptions[*tg.InlineKeyboardMarkup]) (*tg.Message, error) {
	p := params{}
	g.params(p)
	if len(opts) > 0 {
		opts[0].embed(p)
	}
	return chatMethod[*tg.Message](c, "sendGame", p)
}

// StopPoll stops a poll which was sent by the bot.
func (c *Chat) StopPoll(msgID int, replyMarkup ...tg.InlineKeyboardMarkup) (*tg.Poll, error) {
	p := params{}
	p.setInt("message_id", msgID)
	if len(replyMarkup) > 0 {
		p.setJSON("reply_markup", replyMarkup[0])
	}
	return chatMethod[*tg.Poll](c, "stopPoll", p)
}

// DeleteMessage deletes a message, including service messages.
func (c *Chat) DeleteMessage(msgID int) error {
	p := params{}.setInt("message_id", msgID)
	return c.method("deleteMessage", p)
}

// Ban contains parameters for banning a chat member.
type Ban struct {
	UserID         int64
	UntilDate      *int64
	RevokeMessages bool
}

// Ban bans a user in a group, a supergroup or a channel.
func (c *Chat) Ban(b Ban) error {
	p := params{}
	p.setInt64("user_id", b.UserID)
	if b.UntilDate != nil {
		p.setInt64("until_date", *b.UntilDate, true)
	}
	p.setBool("revoke_messages", b.RevokeMessages)
	return c.method("banChatMember", p)
}

// Unban unbans a previously banned user in a supergroup or channel.
func (c *Chat) Unban(userID int64, onlyIfBanned bool) error {
	p := params{}.setInt64("user_id", userID)
	p.setBool("only_if_banned", onlyIfBanned)
	return c.method("unbanChatMember", p)
}

// Restrict restricts a user in a supergroup.
func (c *Chat) Restrict(userID int64, perms tg.ChatPermissions, until *int64) error {
	p := params{}.setInt64("user_id", userID)
	p.setJSON("permissions", perms)
	if until != nil {
		p.setInt64("until_date", *until, true)
	}
	return c.method("restrictChatMember", p)
}

// Promote promotes or demotes a user in a supergroup or a channel.
func (c *Chat) Promote(userID int64, rights tg.ChatAdministratorRights) error {
	p := params{}.setInt64("user_id", userID)
	p.setBool("is_anonymous", rights.IsAnonymous)
	p.setBool("can_manage_chat", rights.CanManageChat)
	p.setBool("can_delete_messages", rights.CanDeleteMessages)
	p.setBool("can_manage_video_chats", rights.CanManageVideoChats)
	p.setBool("can_restrict_members", rights.CanRestrictMembers)
	p.setBool("can_promote_members", rights.CanPromoteMembers)
	p.setBool("can_change_info", rights.CanChangeInfo)
	p.setBool("can_invite_users", rights.CanInviteUsers)
	p.setBool("can_post_messages", rights.CanPostMessages)
	p.setBool("can_edit_messages", rights.CanEditMessages)
	p.setBool("can_pin_messages", rights.CanPinMessages)
	return c.method("promoteChatMember", p)
}

// SetAdminTitle sets a custom title for an administrator in a supergroup promoted by the bot.
func (c *Chat) SetAdminTitle(userID int64, title string) error {
	p := params{}
	p.setInt64("user_id", userID)
	p.set("custom_title", title)
	return c.method("setChatAdministratorCustomTitle", p)
}

// BanSenderChat bans a channel chat in a supergroup or a channel.
func (c *Chat) BanSenderChat(senderID int64) error {
	p := params{}.setInt64("sender_chat_id", senderID)
	return c.method("banChatSenderChat", p)
}

// UnbanSenderChat unbans a previously banned channel chat in a supergroup or channel.
func (c *Chat) UnbanSenderChat(senderID int64) error {
	p := params{}.setInt64("sender_chat_id", senderID)
	return c.method("unbanChatSenderChat", p)
}

// SetPermissions sets default chat permissions for all members.
func (c *Chat) SetPermissions(perms tg.ChatPermissions) error {
	p := params{}.setJSON("permissions", perms)
	return c.method("setChatPermissions", p)
}

// ExportInviteLink generates a new primary invite link for a chat;
// any previously generated primary link is revoked.
func (c *Chat) ExportInviteLink() (string, error) {
	return chatMethod[string](c, "exportChatInviteLink", nil)
}

// InviteLink contains parameters for manipulations with invite links.
type InviteLink struct {
	Name               string
	ExpireDate         int64
	MemberLimit        int
	CreatesJoinRequest bool
}

func (i InviteLink) params(p params) {
	p.set("name", i.Name)
	p.setInt64("expire_date", i.ExpireDate)
	p.setInt("member_limit", i.MemberLimit)
	p.setBool("creates_join_request", i.CreatesJoinRequest)
}

// CreateInviteLink creates an additional invite link for a chat.
func (c *Chat) CreateInviteLink(i InviteLink) (*tg.ChatInviteLink, error) {
	p := params{}
	i.params(p)
	return chatMethod[*tg.ChatInviteLink](c, "createChatInviteLink", p)
}

// EditInviteLink edits a non-primary invite link created by the bot.
func (c *Chat) EditInviteLink(link string, i InviteLink) (*tg.ChatInviteLink, error) {
	p := params{}.set("invite_link", link)
	i.params(p)
	return chatMethod[*tg.ChatInviteLink](c, "editChatInviteLink", p)
}

// RevokeInviteLink revokes an invite link created by the bot.
func (c *Chat) RevokeInviteLink(link string) (*tg.ChatInviteLink, error) {
	p := params{}.set("invite_link", link)
	return chatMethod[*tg.ChatInviteLink](c, "revokeChatInviteLink", p)
}

// ApproveJoinRequest approves a chat join request.
func (c *Chat) ApproveJoinRequest(userID int64) error {
	p := params{}.setInt64("user_id", userID)
	return c.method("approveChatJoinRequest", p)
}

// DeclineJoinRequest declines a chat join request.
func (c *Chat) DeclineJoinRequest(userID int64) error {
	p := params{}.setInt64("user_id", userID)
	return c.method("declineChatJoinRequest", p)
}

// SetPhoto sets a new profile photo for the chat.
func (c *Chat) SetPhoto(photo *tg.InputFile) error {
	return c.method("setChatPhoto", nil, file{
		field:         "photo",
		FileSignature: photo,
	})
}

// DeletePhoto deletes a chat photo.
func (c *Chat) DeletePhoto() error {
	return c.method("deleteChatPhoto", nil)
}

// SetTitle change the title of a chat.
func (c *Chat) SetTitle(title string) error {
	p := params{}.set("title", title)
	return c.method("setChatTitle", p)
}

// SetDescription changes the description of a group, a supergroup or a channel.
func (c *Chat) SetDescription(description string) error {
	p := params{}.set("description", description)
	return c.method("setChatDescription", p)
}

// PinMessage adds a message to the list of pinned messages in a chat.
func (c *Chat) PinMessage(msgID int, notify bool) error {
	p := params{}.setInt("message_id", msgID)
	p.setBool("disable_notification", !notify)
	return c.method("pinChatMessage", p)
}

// UnpinMessage removes a message from the list of pinned messages in a chat.
func (c *Chat) UnpinMessage(msgID int) error {
	p := params{}.setInt("message_id", msgID)
	return c.method("unpinChatMessage", p)
}

// UnpinAllMessages clears the list of pinned messages in a chat.
func (c *Chat) UnpinAllMessages() error {
	return c.method("unpinAllChatMessages", nil)
}

// SetStickerSet sets a new group sticker set for a supergroup.
func (c *Chat) SetStickerSet(stickerSet string) error {
	p := params{}.set("sticker_set_name", stickerSet)
	return c.method("setChatStickerSet", p)
}

// DeleteStickerSet deletes a group sticker set from a supergroup.
func (c *Chat) DeleteStickerSet() error {
	return c.method("deleteChatStickerSet", nil)
}

// SetMenuButton changes the bot's menu button in a private chat.
func (c *Chat) SetMenuButton(menu tg.MenuButton) error {
	p := params{}.setJSON("menu_button", menu)
	return c.method("setChatMenuButton", p)
}

// GetMenuButton returns the current value of the bot's menu button in a private chat.
func (c *Chat) GetMenuButton() (*tg.MenuButton, error) {
	return chatMethod[*tg.MenuButton](c, "getChatMenuButton", nil)
}
