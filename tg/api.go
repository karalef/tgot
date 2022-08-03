package tg

import (
	"encoding/json"
	"fmt"
)

// DefaultAPIURL is a default url for telegram api.
const DefaultAPIURL = "https://api.telegram.org/bot"

// DefaultFileURL is a default url for downloading files.
const DefaultFileURL = "https://api.telegram.org/file/bot"

// APIResponse represents telegram api response.
type APIResponse[T any] struct {
	Ok     bool `json:"ok"`
	Result T    `json:"result"`
	*APIError
}

// APIError describes telegram api error.
type APIError struct {
	Code        int    `json:"error_code"`
	Description string `json:"description"`
	Parameters  *struct {
		MigrateTo  *int64 `json:"migrate_to_chat_id"`
		RetryAfter *int   `json:"retry_after"`
	} `json:"parameters"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("telegram: (%d) %s", e.Code, e.Description)
}

// Update object represents an incoming update.
type Update struct {
	ID                int                `json:"update_id"`
	Message           *Message           `json:"message"`
	EditedMessage     *Message           `json:"edited_message"`
	ChannelPost       *Message           `json:"channel_post"`
	EditedChannelPost *Message           `json:"edited_channel_post"`
	CallbackQuery     *CallbackQuery     `json:"callback_query"`
	InlineQuery       *InlineQuery       `json:"inline_query"`
	InlineChosen      *InlineChosen      `json:"chosen_inline_result"`
	ShippingQuery     *ShippingQuery     `json:"shipping_query"`
	PreCheckoutQuery  *PreCheckoutQuery  `json:"pre_checkout_query"`
	Poll              *Poll              `json:"poll"`
	PollAnswer        *PollAnswer        `json:"poll_answer"`
	MyChatMember      *ChatMemberUpdated `json:"my_chat_member"`
	ChatMember        *ChatMemberUpdated `json:"chat_member"`
	ChatJoinRequest   *ChatJoinRequest   `json:"chat_join_request"`
}

// WebhookInfo describes the current status of a webhook.
type WebhookInfo struct {
	URL                string   `json:"url"`
	HasCustomCert      bool     `json:"has_custom_certificate"`
	PendingUpdateCount int      `json:"pending_update_count"`
	IPAddress          string   `json:"ip_address,omitempty"`
	LastErrorDate      int64    `json:"last_error_date,omitempty"`
	LastErrorMessage   string   `json:"last_error_message,omitempty"`
	LastSyncErrorDate  int64    `json:"last_synchronization_error_date,omitempty"`
	MaxConnections     int      `json:"max_connections,omitempty"`
	AllowedUpdates     []string `json:"allowed_updates,omitempty"`
}

// Command represents a bot command.
type Command struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// CommandScopeType represents scope type.
type CommandScopeType string

// all available command scope types.
const (
	ScopeTypeDefault         = CommandScopeType("default")
	ScopeTypeAllPrivateChats = CommandScopeType("all_private_chats")
	ScopeTypeAllGroupChats   = CommandScopeType("all_group_chats")
	ScopeTypeAllChatAdmins   = CommandScopeType("all_chat_administrators")
	ScopeTypeChat            = CommandScopeType("chat")
	ScopeTypeChatAdmins      = CommandScopeType("chat_administrators")
	ScopeTypeChatMember      = CommandScopeType("chat_member")
)

// CommandScope represents the scope to which bot commands are applied.
type CommandScope struct {
	Type   CommandScopeType `json:"type"`
	ChatID int64            `json:"chat_id,omitempty"`
	UserID int64            `json:"user_id,omitempty"`
}

// MenuButton describes the bot's menu button in a private chat.
type MenuButton struct {
	typ string
	txt string
	app *WebAppInfo
}

// MarshalJSON implements json.Marshaler.
func (b MenuButton) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string      `json:"type"`
		Text   string      `json:"text,omitempty"`
		WebApp *WebAppInfo `json:"web_app,omitempty"`
	}{
		Type:   b.typ,
		Text:   b.txt,
		WebApp: b.app,
	})
}

// MenuButtonCommands represents a menu button, which opens the bot's list of commands.
func MenuButtonCommands() *MenuButton {
	return &MenuButton{
		typ: "commands",
	}
}

// MenuButtonWebApp represents a menu button, which launches a Web App.
func MenuButtonWebApp(text string, webApp *WebAppInfo) *MenuButton {
	return &MenuButton{
		typ: "web_app",
		txt: text,
		app: webApp,
	}
}

// MenuButtonDefault describes that no specific value for the menu button was set.
func MenuButtonDefault() *MenuButton {
	return &MenuButton{
		typ: "default",
	}
}
