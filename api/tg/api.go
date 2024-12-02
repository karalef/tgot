package tg

import (
	"fmt"

	"github.com/karalef/tgot/api/internal/oneof"
)

// Response represents telegram api response.
type Response[Type any] struct {
	Ok     bool `json:"ok"`
	Result Type `json:"result"`
	*Error
}

// Error describes telegram api error.
type Error struct {
	Code        int    `json:"error_code"`
	Description string `json:"description"`
	Parameters  *struct {
		MigrateTo  *int64 `json:"migrate_to_chat_id"`
		RetryAfter *int   `json:"retry_after"`
	} `json:"parameters"`
}

func (e *Error) Error() string {
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
	MenuButtonTypeCommands MenuButtonType = "commands"
	MenuButtonTypeWebApp   MenuButtonType = "web_app"
	MenuButtonTypeDefault  MenuButtonType = "default"
)

var menuButtonTypes = oneof.Map[MenuButtonType]{
	MenuButtonTypeCommands: MenuButtonCommands{},
	MenuButtonTypeWebApp:   MenuButtonWebApp{},
	MenuButtonTypeDefault:  MenuButtonDefault{},
}

type oneOfMenuButton struct{}

func (oneOfMenuButton) New(t MenuButtonType) (oneof.Value[MenuButtonType], bool) {
	return menuButtonTypes.New(t)
}

// MenuButton describes the bot's menu button in a private chat.
type MenuButton = oneof.Object[MenuButtonType, oneOfMenuButton]

// MenuButtonCommands represents a menu button, which opens the bot's list of commands.
type MenuButtonCommands struct{}

func (MenuButtonCommands) Type() MenuButtonType { return MenuButtonTypeCommands }

// MenuButtonWebApp represents a menu button, which launches a Web App.
type MenuButtonWebApp struct {
	Text   string     `json:"text"`
	WebApp WebAppInfo `json:"web_app"`
}

func (MenuButtonWebApp) Type() MenuButtonType { return MenuButtonTypeWebApp }

// MenuButtonDefault describes that no specific value for the menu button was set.
type MenuButtonDefault struct{}

func (MenuButtonDefault) Type() MenuButtonType { return MenuButtonTypeDefault }

// BotName represents the bot's name.
type BotName struct {
	Name string `json:"name"`
}

// BotDescription represents the bot's description.
type BotDescription struct {
	Description string `json:"description"`
}

// BotShortDescription represents the bot's short description.
type BotShortDescription struct {
	ShortDescription string `json:"short_description"`
}

// BusinessConnection describes the connection of the bot with a business account.
type BusinessConnection struct {
	ID         string `json:"id"`
	User       User   `json:"user"`
	UserChatID int64  `json:"user_chat_id"`
	Date       int64  `json:"date"`
	CanReply   bool   `json:"can_reply"`
	IsEnabled  bool   `json:"is_enabled"`
}
