package api

import "github.com/karalef/tgot/tg"

// CommandsData contains parameters for setMyCommands method.
type CommandsData struct {
	Commands []tg.Command
	Scope    tg.CommandScope
	Lang     string
}

func (cd *CommandsData) data() Data {
	d := NewData()
	if cd == nil {
		return d
	}
	d.Set("language_code", cd.Lang)
	d.SetJSON("scope", cd.Scope)
	d.SetJSON("commands", cd.Commands)
	return d
}

// GetCommands returns the current list of the bot's commands for the given scope and user language.
func (a *API) GetCommands(s tg.CommandScope, lang string) ([]tg.Command, error) {
	cd := CommandsData{
		Scope: s,
		Lang:  lang,
	}
	return Request[[]tg.Command](a, "getMyCommands", cd.data())
}

// SetCommands changes the list of the bot's commands.
func (a *API) SetCommands(cd *CommandsData) error {
	return a.Request("setMyCommands", cd.data())
}

// DeleteCommands deletes the list of the bot's commands for the given scope and user language.
func (a *API) DeleteCommands(s tg.CommandScope, lang string) error {
	cd := CommandsData{
		Scope: s,
		Lang:  lang,
	}
	return a.Request("deleteMyCommands", cd.data())
}

// SetDefaultAdminRights changes the default administrator rights requested by the bot
// when it's added as an administrator to groups or channels.
func (a *API) SetDefaultAdminRights(rights *tg.ChatAdministratorRights, forChannels bool) error {
	d := NewData().SetJSON("rights", rights).SetBool("for_channels", forChannels)
	return a.Request("setMyDefaultAdministratorRights", d)
}

// GetDefaultAdminRights returns the current default administrator rights of the bot.
func (a *API) GetDefaultAdminRights(forChannels bool) (*tg.ChatAdministratorRights, error) {
	d := NewData().SetBool("for_channels", forChannels)
	return Request[*tg.ChatAdministratorRights](a, "getMyDefaultAdministratorRights", d)
}

// SetDefaultChatMenuButton changes the bot's default menu button.
//
// This method is a wrapper for setChatMenuButton without specifying the chat id.
func (a *API) SetDefaultChatMenuButton(menu tg.MenuButton) error {
	d := NewData().SetJSON("menu_button", menu)
	return a.Request("setChatMenuButton", d)
}

// GetDefaultChatMenuButton returns the current value of the bot's default menu button.
//
// This method is a wrapper for getChatMenuButton without specifying the chat id.
func (a *API) GetDefaultChatMenuButton() (*tg.MenuButton, error) {
	return Request[*tg.MenuButton](a, "getChatMenuButton")
}
