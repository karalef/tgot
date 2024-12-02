package tgot

import (
	stdcontext "context"
	"errors"
	"io"
	"sync/atomic"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// NewWithToken creates new bit with specified token and default http client.
func NewWithToken(token string, handler Handler) (*Bot, error) {
	a, err := api.New(token, "", "", nil)
	if err != nil {
		return nil, err
	}
	return New(stdcontext.Background(), a, handler)
}

// New creates new bot.
func New(ctx stdcontext.Context, api *api.API, handler Handler) (*Bot, error) {
	if api == nil {
		return nil, errors.New("nil api")
	}
	b := &Bot{
		api: api,
		h:   handler,
	}
	_, err := b.GetMe()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Bot type.
type Bot struct {
	ctx *context
	api *api.API
	err atomic.Pointer[error]
	me  tg.User

	h Handler
}

// Handler represents updates handler.
type Handler interface {
	// Allowed returns list of allowed updates.
	// If the list is nil, all updates are allowed.
	Allowed() []string

	// Handle handles the update.
	// If the handler returns an error, it will be set as bot error,
	// after which the bot will not handle the future updates.
	Handle(Empty, *tg.Update) error
}

// API returns api object.
func (b *Bot) API() *api.API { return b.api }

// Allowed returns list of allowed updates.
// If the list is nil, all updates are allowed.
func (b *Bot) Allowed() []string { return b.h.Allowed() }

// SetError changes the state of the bot to an error.
// It does nothing if the bot is already in an error state.
// The error state prevents any updates handling..
func (b *Bot) SetError(err error) {
	b.err.CompareAndSwap(nil, &err)
}

// Err returns bot error.
func (b *Bot) Err() error {
	if e := b.err.Load(); e != nil {
		return *e
	}
	return nil
}

// Handle calls underlying handler.
func (b *Bot) Handle(upd *tg.Update) error {
	err := b.Err()
	if err != nil {
		return err
	}

	err = b.h.Handle(b.ctx, upd)
	if err != nil {
		b.SetError(err)
	}
	return b.Err()
}

// Me returns current bot as tg.User (cached after GetMe).
// Use it only if you doesn't update the bot info while runtime.
func (b *Bot) Me() tg.User { return b.me }

// GetMe returns basic information about the bot in form of a User object.
func (b *Bot) GetMe() (*tg.User, error) {
	me, err := method[*tg.User](b.ctx, "getMe")
	if err != nil {
		return nil, err
	}
	b.me = *me
	return me, nil
}

// GetFile returns basic information about a file and prepares it for downloading.
func (b *Bot) GetFile(fileID string) (*tg.File, error) {
	d := api.NewData().Set("file_id", fileID)
	return method[*tg.File](b.ctx, "getFile", d)
}

// Download downloads file as io.ReadCloser from Telegram servers.
func (b *Bot) Download(filepath string) (io.ReadCloser, error) {
	return b.api.DownloadFile(b.ctx, filepath)
}

// DownloadFull downloads file from Telegram servers.
func (b *Bot) DownloadFull(filepath string) ([]byte, error) {
	rc, err := b.Download(filepath)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// DownloadFile downloads file as io.ReadCloser from Telegram servers
// by file id.
func (b *Bot) DownloadFile(fid string) (io.ReadCloser, error) {
	f, err := b.GetFile(fid)
	if err != nil {
		return nil, err
	}
	return b.Download(f.FilePath)
}

// DownloadFullFile downloads file from Telegram servers by file id.
func (b *Bot) DownloadFullFile(fid string) ([]byte, error) {
	f, err := b.GetFile(fid)
	if err != nil {
		return nil, err
	}
	return b.DownloadFull(f.FilePath)
}

// SetDefaultAdminRights changes the default administrator rights requested by the bot
// when it's added as an administrator to groups or channels.
func (b *Bot) SetDefaultAdminRights(rights *tg.ChatAdministratorRights, forChannels bool) error {
	d := api.NewData().SetJSON("rights", rights).SetBool("for_channels", forChannels)
	return b.ctx.method("setMyDefaultAdministratorRights", d)
}

// GetDefaultAdminRights returns the current default administrator rights of the bot.
func (b *Bot) GetDefaultAdminRights(forChannels bool) (*tg.ChatAdministratorRights, error) {
	d := api.NewData().SetBool("for_channels", forChannels)
	return method[*tg.ChatAdministratorRights](b.ctx, "getMyDefaultAdministratorRights", d)
}

// SetDefaultChatMenuButton changes the bot's default menu button.
//
// This method is a wrapper for setChatMenuButton without specifying the chat id.
func (b *Bot) SetDefaultChatMenuButton(menu tg.MenuButton) error {
	d := api.NewData().SetJSON("menu_button", menu)
	return b.ctx.method("setChatMenuButton", d)
}

// GetDefaultChatMenuButton returns the current value of the bot's default menu button.
//
// This method is a wrapper for getChatMenuButton without specifying the chat id.
func (b *Bot) GetDefaultChatMenuButton() (*tg.MenuButton, error) {
	return method[*tg.MenuButton](b.ctx, "getChatMenuButton")
}

// SetName changes the bot's name.
func (b *Bot) SetName(name, lang string) error {
	d := api.NewData().Set("name", name)
	d.Set("language_code", lang)
	return b.ctx.method("setMyName", d)
}

// GetName returns the current bot name for the given user language.
func (b *Bot) GetName(lang string) (*tg.BotName, error) {
	d := api.NewData().Set("language_code", lang)
	return method[*tg.BotName](b.ctx, "getMyName", d)
}

// SetDescription changes the bot's description, which is shown in the chat with the bot if the chat is empty.
func (b *Bot) SetDescription(description, lang string) error {
	d := api.NewData().Set("description", description)
	d.Set("language_code", lang)
	return b.ctx.method("setMyDescription", d)
}

// GetDescription returns the current bot description for the given user language.
func (b *Bot) GetDescription(lang string) (*tg.BotDescription, error) {
	d := api.NewData().Set("language_code", lang)
	return method[*tg.BotDescription](b.ctx, "getMyDescription", d)
}

// SetShortDescription changes the bot's short description, which is shown on the bot's profile page and
// is sent together with the link when users share the bot.
func (b *Bot) SetShortDescription(shortDescription, lang string) error {
	d := api.NewData().Set("short_description", shortDescription)
	d.Set("language_code", lang)
	return b.ctx.method("setMyShortDescription", d)
}

// GetShortDescription returns the current bot short description for the given user language.
func (b *Bot) GetShortDescription(lang string) (*tg.BotShortDescription, error) {
	d := api.NewData().Set("language_code", lang)
	return method[*tg.BotShortDescription](b.ctx, "getMyShortDescription", d)
}

// GetBusinessConnection returns information about the connection of the bot with a business account.
func (b *Bot) GetBusinessConnection(id string) (*tg.BusinessConnection, error) {
	d := api.NewData().Set("business_connection_id", id)
	return method[*tg.BusinessConnection](b.ctx, "getBusinessConnection", d)
}

// GetStarTransactions returns the bot's Telegram Star transactions in chronological order.
func (b *Bot) GetStarTransactions(offset uint, limit uint8) (*tg.StarTransactions, error) {
	d := api.NewData().SetInt("offset", int(offset)).SetInt("limit", int(limit))
	return method[*tg.StarTransactions](b.ctx, "getStarTransactions", d)
}

// CommandsData contains parameters for setMyCommands method.
type CommandsData struct {
	Commands []tg.Command    `tg:"commands"`
	Scope    tg.CommandScope `tg:"scope"`
	Lang     string          `tg:"language_code"`
}

// GetCommands returns the current list of the bot's commands for the given scope and user language.
func (b *Bot) GetCommands(s tg.CommandScope, lang string) ([]tg.Command, error) {
	return method[[]tg.Command](b.ctx, "getMyCommands", api.NewDataFrom(CommandsData{
		Scope: s,
		Lang:  lang,
	}))
}

// SetCommands changes the list of the bot's commands.
func (b *Bot) SetCommands(cd *CommandsData) error {
	return b.ctx.method("setMyCommands", api.NewDataFrom(cd))
}

// DeleteCommands deletes the list of the bot's commands for the given scope and user language.
func (b *Bot) DeleteCommands(s tg.CommandScope, lang string) error {
	return b.ctx.method("deleteMyCommands", api.NewDataFrom(CommandsData{
		Scope: s,
		Lang:  lang,
	}))
}

// GetUpdates contains parameters for getUpdates method.
type GetUpdates struct {
	Offset  int      `tg:"offset"`
	Limit   int      `tg:"limit"`
	Timeout int      `tg:"timeout"`
	Allowed []string `tg:"allowed"`
}

// GetUpdates receives incoming updates using long polling.
//
// It is not efficient method because it serializes all parameters on each call
// and uses the bot's context for request.
func (b *Bot) GetUpdates(gu GetUpdates) ([]tg.Update, error) {
	return method[[]tg.Update](b.ctx, "getUpdates", api.NewDataFrom(gu))
}

// WebhookData contains parameters for setWebhook method.
type WebhookData struct {
	URL            string        `tg:"url"`
	Certificate    *tg.InputFile `tg:"certificate"`
	IPAddress      string        `tg:"ip_address"`
	MaxConnections int           `tg:"max_connections"`
	AllowedUpdates []string      `tg:"allowed_updates"`
	DropPending    bool          `tg:"drop_pending_updates"`
	SecretToken    string        `tg:"secret_token"`
}

// SetWebhook specifies a webhook URL.
// Use this method to specify a URL and receive incoming updates via an outgoing webhook.
func (b *Bot) SetWebhook(wd WebhookData) (bool, error) {
	return method[bool](b.ctx, "setWebhook", api.NewDataFrom(wd))
}

// DeleteWebhook removes webhook integration if you decide to switch back to getUpdates.
func (b *Bot) DeleteWebhook(dropPending bool) (bool, error) {
	d := api.NewData().SetBool("drop_pending_updates", dropPending)
	return method[bool](b.ctx, "deleteWebhook", d)
}

// GetWebhookInfo returns current webhook status.
func (b *Bot) GetWebhookInfo() (*tg.WebhookInfo, error) {
	return method[*tg.WebhookInfo](b.ctx, "getWebhookInfo", nil)
}

// CreateInvoiceLink creates a link for an invoice.
func (b *Bot) CreateInvoiceLink(l tg.InputInvoiceMessageContent) (string, error) {
	return method[string](b.ctx, "createInvoiceLink", api.NewDataFrom(l, "json"))
}

// LogOut method.
//
// Use this method to log out from the cloud Bot API server before launching the bot locally.
func (b *Bot) LogOut() error {
	return b.ctx.method("logOut")
}

// Close method.
//
// Use this method to close the bot instance before moving it from one local server to another.
func (b *Bot) Close() error {
	return b.ctx.method("close")
}

// GetCustomEmojiStickers returns information about custom emoji stickers by their identifiers.
func (b *Bot) GetCustomEmojiStickers(ids ...string) ([]tg.Sticker, error) {
	d := api.NewData().SetJSON("custom_emoji_ids", ids)
	return method[[]tg.Sticker](b.ctx, "getCustomEmojiStickers", d)
}

// GetForumTopicIconStickers returns custom emoji stickers,
// which can be used as a forum topic icon by any user.
func (b *Bot) GetForumTopicIconStickers() ([]tg.Sticker, error) {
	return method[[]tg.Sticker](b.ctx, "getForumTopicIconStickers", nil)
}
