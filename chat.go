package bot

import (
	"tghwbot/bot/tg"
)

// Chat represents chat api.
type Chat struct {
	ctx    commonContext
	chatID string
}

// GetInfo returns up to date information about the chat.
func (c *Chat) GetInfo() *tg.Chat {
	return api[*tg.Chat](c.ctx, "getChat", params{
		"chat_id": {c.chatID},
	})
}

// GetAdmins returns a list of administrators in a chat.
func (c *Chat) GetAdmins() []tg.ChatMember {
	return api[[]tg.ChatMember](c.ctx, "getChatAdministrators", params{
		"chat_id": {c.chatID},
	})
}

// MemberCount returns the number of members in a chat.
func (c *Chat) MemberCount() int {
	return api[int](c.ctx, "getChatMemberCount", params{
		"chat_id": {c.chatID},
	})
}

// GetMember returns information about a member of a chat.
func (c *Chat) GetMember(userID int64) *tg.ChatMember {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("user_id", userID)
	return api[*tg.ChatMember](c.ctx, "getChatMember", p)
}

// Leave a group, supergroup or channel.
func (c *Chat) Leave() {
	api[bool](c.ctx, "leaveChat", params{
		"chat_id": {c.chatID},
	})
}

// Forward contains paramenets for forwarding the message.
type Forward struct {
	MessageID           int
	DisableNotification bool
	ProtectContent      bool
}

// Forward forwards messages of any kind.
// Service messages can't be forwarded.
func (c *Chat) Forward(from *Chat, fwd Forward) *tg.Message {
	p := params{}.set("chat_id", c.chatID)
	p.set("from_chat_id", from.chatID)
	p.setInt("message_id", fwd.MessageID)
	p.setBool("disable_notification", fwd.DisableNotification)
	p.setBool("protect_content", fwd.ProtectContent)
	return api[*tg.Message](c.ctx, "forwardMessage", p)
}

// ForwardTo forwards the message to the specified chat instead of current.
func (c *Chat) ForwardTo(to *Chat, fwd Forward) *tg.Message {
	return to.Forward(c, fwd)
}

// Copy contains parameters for copying the message.
type Copy struct {
	MessageID int
	CaptionData
}

// Copy copies messages of any kind.
// Service messages and invoice messages can't be copied.
func (c *Chat) Copy(from *Chat, cp Copy, opts ...SendOptions) *tg.Message {
	p := params{}.set("chat_id", c.chatID)
	p.set("from_chat_id", from.chatID)
	p.setInt("message_id", cp.MessageID)
	cp.CaptionData.embed(p)
	if len(opts) > 0 {
		opts[0].embed(p)
	}
	return api[*tg.Message](c.ctx, "copyMessage", p)
}

// CopyTo copies the message to the specified chat instead of current.
func (c *Chat) CopyTo(to *Chat, cp Copy, opts ...SendOptions) *tg.Message {
	return to.Copy(c, cp)
}

// SendText sends just a text.
func (c *Chat) SendText(text string) *tg.Message {
	return c.Send(NewMessage(text))
}

// Send sends any Sendable object.
func (c *Chat) Send(s Sendable, opts ...SendOptions) *tg.Message {
	if s == nil {
		return nil
	}

	m := "send" + s.what()
	p := params{}.set("chat_id", c.chatID)
	s.params(p)
	if len(opts) > 0 {
		opts[0].embed(p)
	}
	var files []File
	if f, ok := s.(Fileable); ok {
		files = f.files()
	}

	return api[*tg.Message](c.ctx, m, p, files...)
}

// SendMediaGroup sends a group of photos, videos, documents or audios as an album.
func (c *Chat) SendMediaGroup(mg MediaGroup, opts ...MediaGroupSendOptions) []tg.Message {
	p := params{}.set("chat_id", c.chatID)
	files, err := mg.data(p)
	if err != nil {
		closeCtx(c.ctx, err)
	}
	if len(opts) > 0 {
		opts[0].embed(p)
	}
	return api[[]tg.Message](c.ctx, "sendMediaGroup", p, files...)
}

// SendChatAction sends chat action to tell the user that something
// is happening on the bot's side.
func (c *Chat) SendChatAction(act tg.ChatAction) {
	p := params{}.set("chat_id", c.chatID)
	p.set("action", string(act))
	api[bool](c.ctx, "sendChatAction", p)
}

// SendInvoice sends invoice.
func (c *Chat) SendInvoice(i Invoice, opts ...InvoiceSendOptions) *tg.Message {
	p := params{}.set("chat_id", c.chatID)
	i.params(p)
	if len(opts) > 0 {
		opts[0].embed(p)
	}
	return api[*tg.Message](c.ctx, "sendInvoice", p)
}

// EditText contains parameters for editing message text.
type EditText struct {
	Text                  string
	ParseMode             tg.ParseMode
	Entities              []tg.MessageEntity
	DisableWebPagePreview bool
	ReplyMarkup           *tg.InlineKeyboardMarkup
}

func (e EditText) params(p params) {
	p.set("text", e.Text)
	p.set("parse_mode", string(e.ParseMode))
	p.setJSON("entities", e.Entities)
	p.setBool("disable_web_page_preview", e.DisableWebPagePreview)
	p.setJSON("reply_markup", e.ReplyMarkup)
}

// EditMessageText edits text and game messages.
func (c *Chat) EditMessageText(msgID int, e EditText) *tg.Message {
	p := params{}.set("chat_id", c.chatID)
	p.setInt("message_id", msgID)
	e.params(p)
	return api[*tg.Message](c.ctx, "editMessageText", p)
}

// EditCaption contains parameters for editing message caption.
type EditCaption struct {
	CaptionData
	ReplyMarkup *tg.InlineKeyboardMarkup
}

func (e EditCaption) params(p params) {
	e.CaptionData.embed(p)
	p.setJSON("reply_markup", e.ReplyMarkup)
}

// EditMessageCaption edits captions of messages.
func (c *Chat) EditMessageCaption(msgID int, e EditCaption) *tg.Message {
	p := params{}.set("chat_id", c.chatID)
	p.setInt("message_id", msgID)
	e.params(p)
	return api[*tg.Message](c.ctx, "editMessageCaption", p)
}

// EditMessageMedia edits animation, audio, document, photo, or video messages.
func (c *Chat) EditMessageMedia(msgID int, m tg.MediaInputter, replyMarkup ...tg.InlineKeyboardMarkup) *tg.Message {
	p := params{}.set("chat_id", c.chatID)
	p.setInt("message_id", msgID)
	if len(replyMarkup) > 0 {
		p.setJSON("reply_markup", replyMarkup[0])
	}
	files, err := prepareInputMedia(p, false, m)
	if err != nil {
		closeCtx(c.ctx, err)
	}
	return api[*tg.Message](c.ctx, "editMessageMedia", p, files...)
}

// EditMessageReplyMarkup edits only the reply markup of messages.
func (c *Chat) EditMessageReplyMarkup(msgID int, replyMarkup tg.InlineKeyboardMarkup) *tg.Message {
	p := params{}.set("chat_id", c.chatID)
	p.setInt("message_id", msgID)
	p.setJSON("reply_markup", replyMarkup)
	return api[*tg.Message](c.ctx, "editMessageReplyMarkup", p)
}

// StopPoll stops a poll which was sent by the bot.
func (c *Chat) StopPoll(msgID int, replyMarkup ...tg.InlineKeyboardMarkup) *tg.Poll {
	p := params{}.set("chat_id", c.chatID)
	p.setInt("message_id", msgID)
	if len(replyMarkup) > 0 {
		p.setJSON("reply_markup", replyMarkup[0])
	}
	return api[*tg.Poll](c.ctx, "stopPoll", p)
}

// DeleteMessage deletes a message, including service messages.
func (c *Chat) DeleteMessage(msgID int) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("message_id", int64(msgID))
	api[bool](c.ctx, "deleteMessage", p)
}

// LiveLocation contains parameters for editing the live location.
type LiveLocation struct {
	Long               float32
	Lat                float32
	HorizontalAccuracy *float32
	Heading            int
	AlertRadius        int
	ReplyMarkup        *tg.InlineKeyboardMarkup
}

func (l LiveLocation) params(p params) {
	p.setFloat("latitude", l.Lat)
	p.setFloat("longitude", l.Long)
	if l.HorizontalAccuracy != nil {
		p.setFloat("horizontal_accuracy", *l.HorizontalAccuracy)
	}
	p.setInt("heading", l.Heading)
	p.setInt("proximity_alert_radius", l.AlertRadius)
	p.setJSON("reply_markup", l.ReplyMarkup)
}

// EditLiveLocation edits live location messages.
func (c *Chat) EditLiveLocation(msgID int, loc LiveLocation) *tg.Message {
	p := params{}.set("chat_id", c.chatID)
	p.setInt("message_id", msgID)
	loc.params(p)
	return api[*tg.Message](c.ctx, "editMessageLiveLocation", p)
}

// StopLiveLocation stops updating a live location message before live_period expires.
func (c *Chat) StopLiveLocation(msgID int, replyMarkup *tg.InlineKeyboardMarkup) *tg.Message {
	p := params{}.set("chat_id", c.chatID)
	p.setInt("message_id", msgID)
	p.setJSON("reply_markup", replyMarkup)
	return api[*tg.Message](c.ctx, "stopMessageLiveLocation", p)
}

// Ban contains parameters for banning a chat member.
type Ban struct {
	UserID         int64
	UntilDate      *int64
	RevokeMessages bool
}

// Ban bans a user in a group, a supergroup or a channel.
func (c *Chat) Ban(b Ban) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("user_id", b.UserID)
	if b.UntilDate != nil {
		p.setInt64("until_date", *b.UntilDate)
	}
	p.setBool("revoke_messages", b.RevokeMessages)
	api[bool](c.ctx, "banChatMember", p)
}

// Unban unbans a previously banned user in a supergroup or channel.
func (c *Chat) Unban(userID int64, onlyIfBanned bool) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("user_id", userID)
	p.setBool("only_if_banned", onlyIfBanned)
	api[bool](c.ctx, "unbanChatMember", p)
}

// Restrict restricts a user in a supergroup.
func (c *Chat) Restrict(userID int64, perms tg.ChatPermissions, until *int64) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("user_id", userID)
	p.setJSON("permissions", perms)
	if until != nil {
		p.setInt64("until_date", *until)
	}
	api[bool](c.ctx, "restrictChatMember", p)
}

// Promote promotes or demotes a user in a supergroup or a channel.
func (c *Chat) Promote(userID int64, rights tg.ChatAdministratorRights) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("user_id", userID)
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
	api[bool](c.ctx, "promoteChatMember", p)
}

// SetAdminTitle sets a custom title for an administrator in a supergroup promoted by the bot.
func (c *Chat) SetAdminTitle(userID int64, title string) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("user_id", userID)
	p.set("custom_title", title)
	api[bool](c.ctx, "setChatAdministratorCustomTitle", p)
}

// BanSenderChat bans a channel chat in a supergroup or a channel.
func (c *Chat) BanSenderChat(senderID int64) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("sender_chat_id", senderID)
	api[bool](c.ctx, "banChatSenderChat", p)
}

// UnbanSenderChat unbans a previously banned channel chat in a supergroup or channel.
func (c *Chat) UnbanSenderChat(senderID int64) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("sender_chat_id", senderID)
	api[bool](c.ctx, "unbanChatSenderChat", p)
}

// SetPermissions sets default chat permissions for all members.
func (c *Chat) SetPermissions(perms tg.ChatPermissions) {
	p := params{}.set("chat_id", c.chatID)
	p.setJSON("permissions", perms)
	api[bool](c.ctx, "setChatPermissions", p)
}

// ExportInviteLink generates a new primary invite link for a chat;
// any previously generated primary link is revoked.
func (c *Chat) ExportInviteLink() string {
	p := params{}.set("chat_id", c.chatID)
	return api[string](c.ctx, "exportChatInviteLink", p)
}

// InviteLink contains parameters for manipulations with invite links.
type InviteLink struct {
	Name               string
	ExpireDate         int64
	MemberLimit        int
	CreatesJoinRequest bool
}

// CreateInviteLink creates an additional invite link for a chat.
func (c *Chat) CreateInviteLink(i InviteLink) *tg.ChatInviteLink {
	p := params{}.set("chat_id", c.chatID)
	p.set("name", i.Name)
	p.setInt64("expire_date", i.ExpireDate)
	p.setInt("member_limit", i.MemberLimit)
	p.setBool("creates_join_request", i.CreatesJoinRequest)
	return api[*tg.ChatInviteLink](c.ctx, "createChatInviteLink", p)
}

// EditInviteLink edits a non-primary invite link created by the bot.
func (c *Chat) EditInviteLink(link string, i InviteLink) *tg.ChatInviteLink {
	p := params{}.set("chat_id", c.chatID)
	p.set("invite_link", link)
	p.set("name", i.Name)
	p.setInt64("expire_date", i.ExpireDate)
	p.setInt("member_limit", i.MemberLimit)
	p.setBool("creates_join_request", i.CreatesJoinRequest)
	return api[*tg.ChatInviteLink](c.ctx, "editChatInviteLink", p)
}

// RevokeInviteLink revokes an invite link created by the bot.
func (c *Chat) RevokeInviteLink(link string) *tg.ChatInviteLink {
	p := params{}.set("chat_id", c.chatID)
	p.set("invite_link", link)
	return api[*tg.ChatInviteLink](c.ctx, "revokeChatInviteLink", p)
}

// ApproveJoinRequest approves a chat join request.
func (c *Chat) ApproveJoinRequest(userID int64) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("user_id", userID)
	api[bool](c.ctx, "approveChatJoinRequest", p)
}

// DeclineJoinRequest declines a chat join request.
func (c *Chat) DeclineJoinRequest(userID int64) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt64("user_id", userID)
	api[bool](c.ctx, "declineChatJoinRequest", p)
}

// SetPhoto sets a new profile photo for the chat.
func (c *Chat) SetPhoto(photo *tg.InputFile) {
	p := params{}.set("chat_id", c.chatID)
	api[bool](c.ctx, "setChatPhoto", p, File{
		Field:     "photo",
		InputFile: photo,
	})
}

// DeletePhoto deletes a chat photo.
func (c *Chat) DeletePhoto(photo *tg.InputFile) {
	p := params{}.set("chat_id", c.chatID)
	api[bool](c.ctx, "deleteChatPhoto", p)
}

// SetTitle change the title of a chat.
func (c *Chat) SetTitle(title string) {
	p := params{}.set("chat_id", c.chatID)
	p.set("title", title)
	api[bool](c.ctx, "setChatTitle", p)
}

// SetDescription changes the description of a group, a supergroup or a channel.
func (c *Chat) SetDescription(description string) {
	p := params{}.set("chat_id", c.chatID)
	p.set("description", description)
	api[bool](c.ctx, "setChatDescription", p)
}

// PinMessage adds a message to the list of pinned messages in a chat.
func (c *Chat) PinMessage(msgID int, notify bool) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt("message_id", msgID)
	p.setBool("disable_notification", !notify)
	api[bool](c.ctx, "pinChatMessage", p)
}

// UnpinMessage removes a message from the list of pinned messages in a chat.
func (c *Chat) UnpinMessage(msgID int) {
	p := params{}.set("chat_id", c.chatID)
	p.setInt("message_id", msgID)
	api[bool](c.ctx, "unpinChatMessage", p)
}

// UnpinAllMessages clears the list of pinned messages in a chat.
func (c *Chat) UnpinAllMessages() {
	p := params{}.set("chat_id", c.chatID)
	api[bool](c.ctx, "unpinAllChatMessages", p)
}

// SetStickerSet sets a new group sticker set for a supergroup.
func (c *Chat) SetStickerSet(stickerSet string) {
	p := params{}.set("chat_id", c.chatID)
	p.set("sticker_set_name", stickerSet)
	api[bool](c.ctx, "setChatStickerSet", p)
}

// DeleteStickerSet deletes a group sticker set from a supergroup.
func (c *Chat) DeleteStickerSet() {
	p := params{}.set("chat_id", c.chatID)
	api[bool](c.ctx, "deleteChatStickerSet", p)
}
