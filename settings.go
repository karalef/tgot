package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// CommandsData contains parameters for setMyCommands method.
type CommandsData struct {
	Commands []tg.Command
	Scope    tg.CommandScope
	Lang     string
}

func (cd *CommandsData) data() *api.Data {
	d := api.NewData()
	if cd == nil {
		return d
	}
	d.Set("language_code", cd.Lang)
	d.SetJSON("scope", cd.Scope)
	d.SetJSON("commands", cd.Commands)
	return d
}

// GetCommands returns the current list of the bot's commands for the given scope and user language.
func (b *Bot) GetCommands(s tg.CommandScope, lang string) ([]tg.Command, error) {
	cd := CommandsData{
		Scope: s,
		Lang:  lang,
	}
	return api.Request[[]tg.Command](b.api, "getMyCommands", cd.data())
}

// SetCommands changes the list of the bot's commands.
func (b *Bot) SetCommands(cd *CommandsData) error {
	return b.api.Request("setMyCommands", cd.data())
}

// DeleteCommands deletes the list of the bot's commands for the given scope and user language.
func (b *Bot) DeleteCommands(s tg.CommandScope, lang string) error {
	cd := CommandsData{
		Scope: s,
		Lang:  lang,
	}
	return b.api.Request("deleteMyCommands", cd.data())
}

// SetDefaultAdminRights changes the default administrator rights requested by the bot
// when it's added as an administrator to groups or channels.
func (b *Bot) SetDefaultAdminRights(rights *tg.ChatAdministratorRights, forChannels bool) error {
	d := api.NewData().SetJSON("rights", rights).SetBool("for_channels", forChannels)
	return b.api.Request("setMyDefaultAdministratorRights", d)
}

// GetDefaultAdminRights returns the current default administrator rights of the bot.
func (b *Bot) GetDefaultAdminRights(forChannels bool) (*tg.ChatAdministratorRights, error) {
	d := api.NewData().SetBool("for_channels", forChannels)
	return api.Request[*tg.ChatAdministratorRights](b.api, "getMyDefaultAdministratorRights", d)
}

// SetDefaultChatMenuButton changes the bot's default menu button.
//
// This method is a wrapper for setChatMenuButton without specifying the chat id.
func (b *Bot) SetDefaultChatMenuButton(menu tg.MenuButton) error {
	d := api.NewData().SetJSON("menu_button", menu)
	return b.api.Request("setChatMenuButton", d)
}

// GetDefaultChatMenuButton returns the current value of the bot's default menu button.
//
// This method is a wrapper for getChatMenuButton without specifying the chat id.
func (b *Bot) GetDefaultChatMenuButton() (*tg.MenuButton, error) {
	return api.Request[*tg.MenuButton](b.api, "getChatMenuButton", nil)
}