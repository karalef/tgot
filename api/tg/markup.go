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
	Keyboard    [][]KeyboardButton `json:"keyboard"`
	Resize      bool               `json:"resize_keyboard,omitempty"`
	OneTime     bool               `json:"one_time_keyboard,omitempty"`
	Placeholder string             `json:"input_field_placeholder,omitempty"`
	Selective   bool               `json:"selective,omitempty"`
}

// KeyboardButton represents one button of the reply keyboard.
type KeyboardButton struct {
	Text            string          `json:"text"`
	RequestContact  bool            `json:"request_contact,omitempty"`
	RequestLocation bool            `json:"request_location,omitempty"`
	RequestPoll     *ButtonPollType `json:"request_poll,omitempty"`
	WebApp          *WebAppInfo     `json:"web_app,omitmepty"`
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
	Text                string        `json:"text"`
	URL                 string        `json:"url,omitempty"`
	CallbackData        string        `json:"callback_data,omitempty"`
	WebApp              *WebAppInfo   `json:"web_app,omitempty"`
	LoginURL            *LoginURL     `json:"login_url,omitempty"`
	SwitchInline        string        `json:"switch_inline_query,omitempty"`
	SwitchInlineCurrent string        `json:"switch_inline_query_current_chat,omitempty"`
	CallbackGame        *CallbackGame `json:"callback_game,omitempty"`
	Pay                 bool          `json:"pay,omitempty"`
}

// LoginURL represents a parameter of the inline keyboard button
// used to automatically authorize a user.
type LoginURL struct {
	URL          string `json:"url"`
	ForwardText  string `json:"forward_text,omitempty"`
	BotUsername  string `json:"bot_username,omitempty"`
	RequestWrite bool   `json:"request_write_access,omitempty"`
}
