package tg

import "encoding/json"

// ReplyMarkup interface.
type ReplyMarkup interface {
	replyMarkup()
}

func (*ReplyKeyboardMarkup) replyMarkup()  {}
func (*InlineKeyboardMarkup) replyMarkup() {}
func (*ReplyKeyboardRemove) replyMarkup()  {}
func (*ForceReply) replyMarkup()           {}

var (
	_ ReplyMarkup = &ReplyKeyboardMarkup{}
	_ ReplyMarkup = &InlineKeyboardMarkup{}
	_ ReplyMarkup = &ReplyKeyboardRemove{}
	_ ReplyMarkup = &ForceReply{}
)

// ReplyKeyboardMarkup represents a custom keyboard with reply options.
type ReplyKeyboardMarkup struct {
	Keyboard     [][]KeyboardButton `json:"keyboard"`
	IsPersistent bool               `json:"is_persistent,omitempty"`
	Resize       bool               `json:"resize_keyboard,omitempty"`
	OneTime      bool               `json:"one_time_keyboard,omitempty"`
	Placeholder  string             `json:"input_field_placeholder,omitempty"`
	Selective    bool               `json:"selective,omitempty"`
}

// KeyboardButton represents one button of the reply keyboard.
type KeyboardButton struct {
	Text            string                      `json:"text"`
	RequestUsers    *KeyboardButtonRequestUsers `json:"request_users,omitempty"`
	RequestChat     *KeyboardButtonRequestChat  `json:"request_chat,omitempty"`
	RequestContact  bool                        `json:"request_contact,omitempty"`
	RequestLocation bool                        `json:"request_location,omitempty"`
	RequestPoll     *ButtonPollType             `json:"request_poll,omitempty"`
	WebApp          *WebAppInfo                 `json:"web_app,omitmepty"`
}

// KeyboardButtonRequestUsers defines the criteria used to request a suitable user.
// The identifier of the selected user will be shared with the bot when the corresponding button is pressed.
type KeyboardButtonRequestUsers struct {
	RequestID       int   `json:"request_id"`
	UserIsBot       *bool `json:"user_is_bot,omitempty"`
	UserIsPremium   *bool `json:"user_is_premium,omitempty"`
	MaxQuantity     int   `json:"max_quantity,omitempty"`
	RequestName     bool  `json:"request_name,omitempty"`
	RequestUsername bool  `json:"request_username,omitempty"`
	RequestPhoto    bool  `json:"request_photo,omitempty"`
}

// KeyboardButtonRequestChat defines the criteria used to request a suitable chat.
// The identifier of the selected chat will be shared with the bot when the corresponding button is pressed.
type KeyboardButtonRequestChat struct {
	RequestID       int                      `json:"request_id"`
	ChatIsChannel   bool                     `json:"chat_is_channel"`
	ChatIsForum     *bool                    `json:"chat_is_forum,omitempty"`
	ChatHasUsername *bool                    `json:"chat_has_username,omitempty"`
	ChatIsCreated   bool                     `json:"chat_is_created,omitempty"`
	UserAdminRights *ChatAdministratorRights `json:"user_admin_rights,omitempty"`
	BotAdminRights  *ChatAdministratorRights `json:"bot_admin_rights,omitempty"`
	BotIsMember     bool                     `json:"bot_is_member,omitempty"`
	RequestTitle    bool                     `json:"request_title,omitempty"`
	RequestUsername bool                     `json:"request_username,omitempty"`
	RequestPhoto    bool                     `json:"request_photo,omitempty"`
}

// ButtonPollType represents type of a poll, which is allowed
// to be created and sent when the corresponding button is pressed.
type ButtonPollType struct {
	Type PollType `json:"type,omitempty"`
}

// ReplyKeyboardRemove represents an object, on receipt of which Telegram clients
// will remove the current custom keyboard and display the default letter-keyboard.
type ReplyKeyboardRemove struct {
	Selective bool
}

// MarshalJSON is json.Marshaler implementation.
func (r ReplyKeyboardRemove) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Remove    bool `json:"remove_keyboard"`
		Selective bool `json:"selective,omitempty"`
	}{true, r.Selective})
}

// ForceReply represents an object, on receipt of which Telegram clients
// will display a reply interface to the user (act as if the user has
// selected the bot's message and tapped 'Reply').
type ForceReply struct {
	Placeholder string
	Selective   bool
}

// MarshalJSON is json.Marshaler implementation.
func (f ForceReply) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ForceReply  bool   `json:"force_reply"`
		Placeholder string `json:"input_field_placeholder,omitempty"`
		Selective   bool   `json:"selective,omitempty"`
	}{true, f.Placeholder, f.Selective})
}

// InlineKeyboardMarkup represents an inline keyboard that
// appears right next to the message it belongs to.
type InlineKeyboardMarkup struct {
	Keyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton represents one button of an inline keyboard.
type InlineKeyboardButton struct {
	Text                string                       `json:"text"`
	URL                 string                       `json:"url,omitempty"`
	CallbackData        string                       `json:"callback_data,omitempty"`
	WebApp              *WebAppInfo                  `json:"web_app,omitempty"`
	LoginURL            *LoginURL                    `json:"login_url,omitempty"`
	SwitchInline        string                       `json:"switch_inline_query,omitempty"`
	SwitchInlineCurrent string                       `json:"switch_inline_query_current_chat,omitempty"`
	SwitchInlineChosen  *SwitchInlineQueryChosenChat `json:"switch_inline_query_chosen_chat,omitempty"`
	CallbackGame        *CallbackGame                `json:"callback_game,omitempty"`
	Pay                 bool                         `json:"pay,omitempty"`
}

// SwitchInlineQueryChosenChat represents an inline button that switches the current user to inline mode in a chosen chat, with an optional default inline query.
type SwitchInlineQueryChosenChat struct {
	Query            string `json:"query"`
	AllowUserChat    bool   `json:"allow_user_chat,omitempty"`
	AllowBotChat     bool   `json:"allow_bot_chat,omitempty"`
	AllowGroupChat   bool   `json:"allow_group_chat,omitempty"`
	AllowChannelChat bool   `json:"allow_channel_chat,omitempty"`
}

// LoginURL represents a parameter of the inline keyboard button
// used to automatically authorize a user.
type LoginURL struct {
	URL          string `json:"url"`
	ForwardText  string `json:"forward_text,omitempty"`
	BotUsername  string `json:"bot_username,omitempty"`
	RequestWrite bool   `json:"request_write_access,omitempty"`
}
