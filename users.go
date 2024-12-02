package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// WithUser creates a new User from context and user id.
// It copies the context but resets the context parameters.
func WithUser(ctx BaseContext, userID int64) *User {
	return &User{
		context: ctx.ctx().with(api.NewData().SetInt64("user_id", userID)),
		id:      userID,
	}
}

var _ Context[*User] = &User{}

// Users contains methods for working with user.
type User struct {
	*context
	id int64
}

func (u *User) WithName(name string) *User {
	return &User{
		context: u.context.child(name),
		id:      u.id,
	}
}

// WithChat returns ChatMember with provided chat id.
func (u *User) WithChat(chatID ChatID) *ChatMember {
	return WithChatMember(u, chatID, u.id)
}

// GetPhotos returns a list of profile pictures for a user.
func (u *User) GetPhotos() (*tg.UserProfilePhotos, error) {
	return method[*tg.UserProfilePhotos](u, "getUserProfilePhotos")
}

// RefundStarPayment refunds a successful payment in Telegram Stars.
func (u *User) RefundStarPayment(chargeID string) error {
	d := api.NewData().Set("telegram_payment_charge_id", chargeID)
	return u.method("refundStarPayment", d)
}

// UploadStickerFile uploads a .PNG file with a sticker for later use
// in createNewStickerSet and addStickerToSet methods.
func (u *User) UploadStickerFile(sticker *tg.InputFile, format tg.StickerFormat) (*tg.File, error) {
	d := api.NewData().AddFile("sticker", sticker).Set("sticker_format", string(format))
	return method[*tg.File](u, "uploadStickerFile", d)
}

// WithChatMember returns ChatMember with provided chat id and user id.
// It copies the context but resets the context parameters.
func WithChatMember(ctx BaseContext, chatID ChatID, userID int64) *ChatMember {
	d := api.NewData().SetInt64("user_id", userID)
	chatID.setChatID(d)
	return &ChatMember{context: ctx.ctx().with(d), chat: chatID, user: userID}
}

// ChatMember provides methods for working with chat members.
type ChatMember struct {
	*context
	user int64
	chat ChatID
}

func (m *ChatMember) WithName(name string) *ChatMember {
	return &ChatMember{
		context: m.context.child(name),
		user:    m.user,
		chat:    m.chat,
	}
}

// Chat returns the Chat of the ChatMember.
func (m *ChatMember) Chat() *Chat { return WithChatID(m, m.chat) }

// User returns the User of the ChatMember.
func (m *ChatMember) User() *User { return WithUser(m, m.user) }

// GetUserBoosts returns the list of boosts added to a chat by a user.
func (m *ChatMember) GetBoosts() (*tg.UserChatBoosts, error) {
	return method[*tg.UserChatBoosts](m, "getUserChatBoosts")
}

// Get returns information about a member of a chat.
func (m *ChatMember) Get() (*tg.ChatMember, error) {
	return method[*tg.ChatMember](m, "getChatMember")
}

// ApproveJoinRequest approves a chat join request.
func (m *ChatMember) ApproveJoinRequest() error {
	return m.method("approveChatJoinRequest")
}

// DeclineJoinRequest declines a chat join request.
func (m *ChatMember) DeclineJoinRequest() error {
	return m.method("declineChatJoinRequest")
}

// SetTitle sets a custom title for an administrator in a supergroup promoted by the bot.
func (m *ChatMember) SetTitle(title string) error {
	d := api.NewData().Set("custom_title", title)
	return m.method("setChatAdministratorCustomTitle", d)
}

// RestrictChatMember contains parameters for restrictChatMember method.
type RestrictChatMember struct {
	Permissions            tg.ChatPermissions `tg:"permissions"`
	IndependentPermissions bool               `tg:"use_independent_chat_permissions"`
	Until                  int64              `tg:"until_date"`
}

// Restrict restricts a user in a supergroup.
func (m *ChatMember) Restrict(r RestrictChatMember) error {
	return m.method("restrictChatMember", api.NewDataFrom(r))
}

// Promote promotes or demotes a user in a supergroup or a channel.
func (m *ChatMember) Promote(rights tg.ChatAdministratorRights) error {
	d := api.NewData()
	api.MarshalTo(d, rights, "json")
	return m.method("promoteChatMember", d)
}

// Ban bans a user in a group, a supergroup or a channel.
func (m *ChatMember) Ban(revokeMessages bool, until ...int64) error {
	d := api.NewData().SetBool("revoke_messages", revokeMessages)
	if len(until) > 0 {
		d.SetInt64("until_date", until[0], true)
	}
	return m.method("banChatMember", d)
}

// Unban unbans a previously banned user in a supergroup or channel.
func (m *ChatMember) Unban(onlyIfBanned bool) error {
	d := api.NewData().SetBool("only_if_banned", onlyIfBanned)
	return m.method("unbanChatMember", d)
}
