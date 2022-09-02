package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"tghwbot/bot/internal"
	"tghwbot/bot/tg"
)

// Error type.
type Error struct {
	Method string
	Params params
	Err    tg.APIError
	cancel bool
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s\n%s %v", e.Err.Error(), e.Method, e.Params)
}

func (e *Error) Unwrap() error {
	return e.Err
}

type jsonError struct {
	Method   string
	Params   params
	Response []byte
	Err      error
}

func (e *jsonError) Error() string {
	return fmt.Sprintf("%s\n%s %v\nresponse:\n%s", e.Err.Error(), e.Method, e.Params, string(e.Response))
}

func (e *jsonError) Unwrap() error {
	return e.Err
}

type params url.Values

func (p params) forEach(f func(k, v string) error) error {
	for k, v := range p {
		if err := f(k, v[0]); err != nil {
			return err
		}
	}
	return nil
}

func (p params) set(k, v string) params {
	if v != "" {
		url.Values(p).Set(k, v)
	}
	return p
}

func (p params) setInt(key string, v int) {
	if v != 0 {
		p.set(key, strconv.Itoa(v))
	}
}

func (p params) setInt64(key string, v int64) {
	if v != 0 {
		p.set(key, strconv.FormatInt(v, 10))
	}
}

func (p params) setFloat(key string, v float32) {
	if v != 0 {
		p.set(key, strconv.FormatFloat(float64(v), 'f', 6, 32))
	}
}

func (p params) setBool(key string, v bool) {
	if v {
		p.set(key, strconv.FormatBool(v))
	}
}

func (p params) setJSON(key string, v interface{}) {
	if v != nil && !reflect.ValueOf(v).IsZero() {
		b, _ := json.Marshal(v)
		p.set(key, string(b))
	}
}

type file struct {
	field string
	tg.FileSignature
}

func performRequestContext[T any](ctx context.Context, b *Bot, method string, p params, files ...file) (T, error) {
	var ctype = "application/x-www-form-urlencoded"
	var data io.Reader
	upload := false
	for i := range files {
		if _, r := files[i].FileData(); r != nil {
			upload = true
			break
		}
	}

	if !upload {
		for i := range files {
			urlid, _ := files[i].FileData()
			p.set(files[i].field, urlid)
		}
		data = strings.NewReader(url.Values(p).Encode())
	} else {
		ctype, data = writeMultipart(p, files)
	}

	var nilResult T
	u := b.apiURL + b.token + "/" + method
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, data)
	if err != nil {
		return nilResult, err
	}
	req.Header.Set("Content-Type", ctype)
	resp, err := b.client.Do(req)
	if err != nil {
		switch e := errors.Unwrap(err); e {
		case context.Canceled, context.DeadlineExceeded:
			return nilResult, e
		default:
			return nilResult, err
		}
	}
	defer resp.Body.Close()

	r, raw, err := internal.DecodeJSON[tg.APIResponse[T]](resp.Body)
	if err != nil {
		return nilResult, &jsonError{
			Method:   method,
			Params:   p,
			Response: raw,
			Err:      err,
		}
	}
	if r.APIError != nil {
		e := &Error{
			Method: method,
			Params: p,
			Err:    *r.APIError,
		}
		if e.Err.Code == http.StatusUnauthorized {
			b.cancel(e.Err)
			e.cancel = true
		}
		return nilResult, e
	}

	return r.Result, nil
}

func performRequest[T any](b *Bot, method string, p params, files ...file) (T, error) {
	return performRequestContext[T](context.Background(), b, method, p, files...)
}

func performRequestEmpty(b *Bot, method string, p params) error {
	_, err := performRequest[internal.Empty](b, method, p)
	return err
}

func writeMultipart(p params, files []file) (string, io.Reader) {
	r, w := io.Pipe()
	mp := multipart.NewWriter(w)
	go func() {
		defer func() {
			w.CloseWithError(mp.Close())
		}()

		err := p.forEach(mp.WriteField)
		if err != nil {
			w.CloseWithError(err)
			return
		}

		for _, file := range files {
			if _, ok := file.FileSignature.(*tg.InputFile); !ok {
				urlid, _ := file.FileData()
				err := mp.WriteField(file.field, urlid)
				if err != nil {
					w.CloseWithError(err)
					return
				}
				continue
			}

			name, reader := file.FileData()
			part, err := mp.CreateFormFile(file.field, name)
			if err != nil {
				w.CloseWithError(err)
				return
			}
			if reader == nil {
				continue
			}
			_, err = io.Copy(part, reader)
			if err != nil {
				w.CloseWithError(err)
				return
			}
		}
	}()
	return mp.FormDataContentType(), r
}

func (b *Bot) downloadFile(path string) (io.ReadCloser, error) {
	resp, err := b.client.Get(b.fileURL + b.token + "/" + path)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// GetMe returns basic information about the bot in form of a User object.
func (b *Bot) GetMe() (*tg.User, error) {
	return performRequest[*tg.User](b, "getMe", nil)
}

// LogOut method.
//
// Use this method to log out from the cloud Bot API server before launching the bot locally.
func (b *Bot) LogOut() error {
	return performRequestEmpty(b, "logOut", nil)
}

// Close method.
//
// Use this method to close the bot instance before moving it from one local server to another.
func (b *Bot) Close() error {
	return performRequestEmpty(b, "close", nil)
}

func (b *Bot) getUpdates(ctx context.Context, offset, timeout, limit int, allowed []string) ([]tg.Update, error) {
	p := params{}
	p.setInt("offset", offset)
	p.setInt("limit", limit)
	p.setInt("timeout", timeout)
	p.setJSON("allowed_updates", allowed)

	return performRequestContext[[]tg.Update](ctx, b, "getUpdates", p)
}

type commandParams struct {
	Commands []tg.Command
	Scope    tg.CommandScope
	Lang     string
}

func (p *commandParams) params() params {
	if p == nil {
		return nil
	}
	v := params{}
	v.set("language_code", p.Lang)
	v.setJSON("scope", p.Scope)
	v.setJSON("commands", p.Commands)
	return v
}

func (b *Bot) getCommands(s tg.CommandScope, lang string) ([]tg.Command, error) {
	p := commandParams{
		Scope: s,
		Lang:  lang,
	}
	return performRequest[[]tg.Command](b, "getMyCommands", p.params())
}

func (b *Bot) setCommands(p *commandParams) error {
	return performRequestEmpty(b, "setMyCommands", p.params())
}

func (b *Bot) deleteCommands(s tg.CommandScope, lang string) error {
	p := commandParams{
		Scope: s,
		Lang:  lang,
	}
	return performRequestEmpty(b, "deleteMyCommands", p.params())
}

// SetDefaultAdminRights changes the default administrator rights requested by the bot
// when it's added as an administrator to groups or channels.
func (b *Bot) SetDefaultAdminRights(rights *tg.ChatAdministratorRights, forChannels bool) error {
	p := params{}
	p.setJSON("rights", rights)
	p.setBool("for_channels", forChannels)
	return performRequestEmpty(b, "setMyDefaultAdministratorRights", p)
}

// GetDefaultAdminRights returns the current default administrator rights of the bot.
func (b *Bot) GetDefaultAdminRights(forChannels bool) (*tg.ChatAdministratorRights, error) {
	p := params{}
	p.setBool("for_channels", forChannels)
	return performRequest[*tg.ChatAdministratorRights](b, "getMyDefaultAdministratorRights", p)
}

// SetDefaultChatMenuButton changes the bot's default menu button.
//
// This method is a wrapper for setChatMenuButton without specifying the chat id.
// Full implementation of this method is available in the [Chat].
func (b *Bot) SetDefaultChatMenuButton(menu tg.MenuButton) error {
	p := params{}
	p.setJSON("menu_button", menu)
	return performRequestEmpty(b, "setChatMenuButton", p)
}

// GetDefaultChatMenuButton returns the current value of the bot's default menu button.
//
// This method is a wrapper for getChatMenuButton without specifying the chat id.
// Full implementation of this method is available in the [Chat].
func (b *Bot) GetDefaultChatMenuButton() (*tg.MenuButton, error) {
	return performRequest[*tg.MenuButton](b, "getChatMenuButton", nil)
}
