package tgot

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

	"github.com/karalef/tgot/internal"
	"github.com/karalef/tgot/tg"
)

func makeError[T error](method string, p params, f []file, err T) baseError[T] {
	return baseError[T]{
		Method: method,
		Params: p,
		Files:  f,
		Err:    err,
	}
}

type baseError[T error] struct {
	Method string
	Params params
	Files  []file
	Err    T
}

func (e baseError[T]) Error() string {
	return fmt.Sprintf("%s\n%s %s", e.Err.Error(), e.Method, e.formatData())
}

func (e baseError[T]) Unwrap() error {
	return e.Err
}

func (e baseError[T]) formatData() string {
	if len(e.Files) == 0 {
		return url.Values(e.Params).Encode()
	}
	var sb strings.Builder
	sb.WriteString(url.Values(e.Params).Encode())
	for _, f := range e.Files {
		sb.WriteByte('\n')
		name, r := f.FileData()
		sb.WriteString("(file) " + f.field + ": " + name)
		if r != nil {
			sb.WriteString(" (upload data)")
		}
	}
	return sb.String()
}

// Error represents a telegram api error and also contains method and data.
type Error struct {
	baseError[tg.APIError]
}

// Is implements errors.Is interface.
func (e Error) Is(err error) bool {
	if tge, ok := err.(tg.APIError); ok {
		return e.Err.Code != tge.Code
	}
	return false
}

// HTTPError represents http error.
type HTTPError struct {
	baseError[error]
}

// JSONError represents JSON error.
type JSONError struct {
	baseError[error]
	Response []byte
}

func (e *JSONError) Error() string {
	return e.baseError.Error() + "\nresponse:\n" + string(e.Response)
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

func (p params) setInt(key string, v int, force ...bool) params {
	if v != 0 || len(force) > 0 && force[0] {
		p.set(key, strconv.Itoa(v))
	}
	return p
}

func (p params) setInt64(key string, v int64, force ...bool) params {
	if v != 0 || len(force) > 0 && force[0] {
		p.set(key, strconv.FormatInt(v, 10))
	}
	return p
}

func (p params) setFloat(key string, v float32, force ...bool) params {
	if v != 0 || len(force) > 0 && force[0] {
		p.set(key, strconv.FormatFloat(float64(v), 'f', 6, 32))
	}
	return p
}

func (p params) setBool(key string, v bool) params {
	if v {
		p.set(key, strconv.FormatBool(v))
	}
	return p
}

func (p params) setJSON(key string, v interface{}) params {
	if v != nil && !reflect.ValueOf(v).IsZero() {
		b, _ := json.Marshal(v)
		p.set(key, string(b))
	}
	return p
}

type file struct {
	field string
	tg.FileSignature
}

func (b *Bot) requestEmpty(method string, p params, files ...file) error {
	_, err := request[internal.Empty](b, method, p, files...)
	return err
}

func request[T any](b *Bot, method string, p params, files ...file) (T, error) {
	return requestContext[T](context.Background(), b, method, p, files...)
}

func requestContext[T any](ctx context.Context, b *Bot, method string, p params, files ...file) (T, error) {
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
		return nilResult, &HTTPError{makeError(method, p, files, err)}
	}
	req.Header.Set("Content-Type", ctype)
	resp, err := b.client.Do(req)
	if err != nil {
		switch e := errors.Unwrap(err); e {
		case context.Canceled, context.DeadlineExceeded:
			return nilResult, e
		default:
			return nilResult, &HTTPError{makeError(method, p, files, err)}
		}
	}
	defer resp.Body.Close()

	r, raw, err := internal.DecodeJSON[tg.APIResponse[T]](resp.Body)
	if err != nil {
		return nilResult, &JSONError{
			baseError: makeError(method, p, files, err),
			Response:  raw,
		}
	}
	if r.APIError != nil {
		return nilResult, &Error{makeError(method, p, files, *r.APIError)}
	}

	return r.Result, nil
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
	return request[*tg.User](b, "getMe", nil)
}

// LogOut method.
//
// Use this method to log out from the cloud Bot API server before launching the bot locally.
func (b *Bot) LogOut() error {
	return b.requestEmpty("logOut", nil)
}

// Close method.
//
// Use this method to close the bot instance before moving it from one local server to another.
func (b *Bot) Close() error {
	return b.requestEmpty("close", nil)
}

func (b *Bot) getUpdates(ctx context.Context, offset, timeout, limit int, allowed []string) ([]tg.Update, error) {
	p := params{}
	p.setInt("offset", offset)
	p.setInt("limit", limit)
	p.setInt("timeout", timeout)
	p.setJSON("allowed_updates", allowed)

	return requestContext[[]tg.Update](ctx, b, "getUpdates", p)
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
	return request[[]tg.Command](b, "getMyCommands", p.params())
}

func (b *Bot) setCommands(p *commandParams) error {
	return b.requestEmpty("setMyCommands", p.params())
}

func (b *Bot) deleteCommands(s tg.CommandScope, lang string) error {
	p := commandParams{
		Scope: s,
		Lang:  lang,
	}
	return b.requestEmpty("deleteMyCommands", p.params())
}

// SetDefaultAdminRights changes the default administrator rights requested by the bot
// when it's added as an administrator to groups or channels.
func (b *Bot) SetDefaultAdminRights(rights *tg.ChatAdministratorRights, forChannels bool) error {
	p := params{}
	p.setJSON("rights", rights)
	p.setBool("for_channels", forChannels)
	return b.requestEmpty("setMyDefaultAdministratorRights", p)
}

// GetDefaultAdminRights returns the current default administrator rights of the bot.
func (b *Bot) GetDefaultAdminRights(forChannels bool) (*tg.ChatAdministratorRights, error) {
	p := params{}.setBool("for_channels", forChannels)
	return request[*tg.ChatAdministratorRights](b, "getMyDefaultAdministratorRights", p)
}

// SetDefaultChatMenuButton changes the bot's default menu button.
//
// This method is a wrapper for setChatMenuButton without specifying the chat id.
// Full implementation of this method is available in the [Chat].
func (b *Bot) SetDefaultChatMenuButton(menu tg.MenuButton) error {
	p := params{}.setJSON("menu_button", menu)
	return b.requestEmpty("setChatMenuButton", p)
}

// GetDefaultChatMenuButton returns the current value of the bot's default menu button.
//
// This method is a wrapper for getChatMenuButton without specifying the chat id.
// Full implementation of this method is available in the [Chat].
func (b *Bot) GetDefaultChatMenuButton() (*tg.MenuButton, error) {
	return request[*tg.MenuButton](b, "getChatMenuButton", nil)
}
