package tg

import (
	"fmt"
)

// DefaultAPIURL is a default url for telegram api.
const DefaultAPIURL = "https://api.telegram.org/bot"

// DefaultFileURL is a default url for downloading files.
const DefaultFileURL = "https://api.telegram.org/file/bot"

// APIResponse represents telegram api response.
type APIResponse[Type any] struct {
	Ok     bool `json:"ok"`
	Result Type `json:"result"`
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

func (e APIError) Error() string {
	s := fmt.Sprintf("telegram (%d) %s", e.Code, e.Description)
	if e.Parameters != nil {
		if e.Parameters.MigrateTo != nil {
			s += fmt.Sprintf(" (migrate to %d)", *e.Parameters.MigrateTo)
		}
		if e.Parameters.RetryAfter != nil {
			s += fmt.Sprintf(" (retry after %d)", *e.Parameters.RetryAfter)
		}
	}
	return s
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

// CallbackQuery represents an incoming callback query from a callback
// button in an inline keyboard.
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageID string   `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
	GameShortName   string   `json:"game_short_name"`
}

// WebhookInfo describes the current status of a webhook.
type WebhookInfo struct {
	URL                string   `json:"url"`
	HasCustomCert      bool     `json:"has_custom_certificate"`
	PendingUpdateCount int      `json:"pending_update_count"`
	IPAddress          string   `json:"ip_address"`
	LastErrorDate      int64    `json:"last_error_date"`
	LastErrorMessage   string   `json:"last_error_message"`
	LastSyncErrorDate  int64    `json:"last_synchronization_error_date"`
	MaxConnections     int      `json:"max_connections"`
	AllowedUpdates     []string `json:"allowed_updates"`
}

// ChatID represents chat id.
type ChatID interface {
	string | int64
}

// Command represents a bot command.
type Command struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// CommandScope represents the scope to which bot commands are applied.
type CommandScope interface {
	commandScope()
}

type commandScope[ChatIDType ChatID] struct {
	Type   string     `json:"type"`
	ChatID ChatIDType `json:"chat_id,omitempty"`
	UserID int64      `json:"user_id,omitempty"`
}

func (commandScope[ChatIDType]) commandScope() {}

// CommandScopeDefault returns the default scope of bot commands.
// Default commands are used if no commands with a narrower scope
// are specified for the user.
func CommandScopeDefault() CommandScope {
	return commandScope[int64]{Type: "default"}
}

// CommandScopeAllPrivateChats returns the scope of bot commands,
// covering all private chats.
func CommandScopeAllPrivateChats() CommandScope {
	return commandScope[int64]{Type: "all_private_chats"}
}

// CommandScopeAllGroupChats returns the scope of bot commands,
// covering all group and supergroup chats.
func CommandScopeAllGroupChats() CommandScope {
	return commandScope[int64]{Type: "all_group_chats"}
}

// CommandScopeAllChatAdmins returns the scope of bot commands,
// covering all group and supergroup chat administrators.
func CommandScopeAllChatAdmins() CommandScope {
	return commandScope[int64]{Type: "all_chat_administrators"}
}

// CommandScopeChat returns the scope of bot commands,
// covering a specific chat.
func CommandScopeChat[T ChatID](chatID T) CommandScope {
	return commandScope[T]{Type: "chat", ChatID: chatID}
}

// CommandScopeChatAdmins returns the scope of bot commands,
// covering all administrators of a specific group or supergroup chat.
func CommandScopeChatAdmins[T ChatID](chatID T) CommandScope {
	return commandScope[T]{Type: "chat_administrators", ChatID: chatID}
}

// CommandScopeChatMember returns the scope of bot commands,
// covering a specific member of a group or supergroup chat.
func CommandScopeChatMember[T ChatID](chatID T, userID int64) CommandScope {
	return commandScope[T]{Type: "chat_member", ChatID: chatID, UserID: userID}
}

// MenuButtonType represents menu button type.
type MenuButtonType string

// all available menu button types.
const (
	MenuButtonTypeCommands = MenuButtonType("commands")
	MenuButtonTypeWebApp   = MenuButtonType("web_app")
	MenuButtonTypeDefault  = MenuButtonType("default")
)

// MenuButton describes the bot's menu button in a private chat.
type MenuButton struct {
	Type   MenuButtonType `json:"type"`
	Text   string         `json:"text,omitempty"`
	WebApp *WebAppInfo    `json:"web_app,omitempty"`
}

// MenuButtonCommands represents a menu button, which opens the bot's list of commands.
func MenuButtonCommands() *MenuButton {
	return &MenuButton{Type: MenuButtonTypeCommands}
}

// MenuButtonWebApp represents a menu button, which launches a Web App.
func MenuButtonWebApp(text string, webApp *WebAppInfo) *MenuButton {
	return &MenuButton{
		Type:   MenuButtonTypeWebApp,
		Text:   text,
		WebApp: webApp,
	}
}

// MenuButtonDefault describes that no specific value for the menu button was set.
func MenuButtonDefault() *MenuButton {
	return &MenuButton{Type: MenuButtonTypeDefault}
}
