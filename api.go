package bot

import (
	"context"
	"encoding/json"
	"errors"
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
	Err      error
	Method   string
	Response []byte
}

func (e *Error) Error() string {
	return e.Err.Error()
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

func performRequest[T any](b *Bot, method string, p params, files ...File) (T, error) {
	var ctype = "application/x-www-form-urlencoded"
	var data io.Reader
	if files == nil {
		data = strings.NewReader(url.Values(p).Encode())
	} else {
		ctype, data = writeMultipart(p, files)
	}

	var nilResult T
	u := b.apiURL + b.token + "/" + method
	req, err := http.NewRequestWithContext(b.ctx, http.MethodPost, u, data)
	if err != nil {
		return nilResult, err
	}
	req.Header.Set("Content-Type", ctype)
	resp, err := b.client.Do(req)
	switch e := errors.Unwrap(err); e {
	case context.Canceled, context.DeadlineExceeded:
		return nilResult, e
	default:
		if err != nil {
			return nilResult, err
		}
	}
	defer resp.Body.Close()

	r, raw, err := internal.DecodeJSON[tg.APIResponse[T]](resp.Body)
	if err != nil {
		return nilResult, &Error{
			Err:      err,
			Method:   method,
			Response: raw,
		}
	}
	if r.APIError != nil {
		return nilResult, &Error{
			Err:    r.APIError,
			Method: method,
		}
	}

	return r.Result, nil
}

func performRequestEmpty(b *Bot, method string, p params) error {
	_, err := performRequest[internal.Empty](b, method, p)
	return err
}

func writeMultipart(p params, files []File) (string, io.Reader) {
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
			if d := file.Data(); d != "" {
				err := mp.WriteField(file.Field, d)
				if err != nil {
					w.CloseWithError(err)
					return
				}
				continue
			}

			name, reader := file.UploadData()
			part, err := mp.CreateFormFile(file.Field, name)
			if err != nil {
				w.CloseWithError(err)
				return
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

// File contains field and file data.
type File struct {
	Field string
	*tg.InputFile
}

func (b *Bot) downloadFile(path string) (io.ReadCloser, error) {
	resp, err := b.client.Get(b.fileURL + b.token + "/" + path)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (b *Bot) getMe() (*tg.User, error) {
	return performRequest[*tg.User](b, "getMe", nil)
}

func (b *Bot) logOut() (bool, error) {
	return performRequest[bool](b, "logOut", nil)
}

func (b *Bot) close() (bool, error) {
	return performRequest[bool](b, "close", nil)
}

func (b *Bot) getUpdates(offset, timeout int, allowed ...string) ([]tg.Update, error) {
	p := params{}
	p.setInt("offset", offset)
	//p.setInt("limit", limit)
	p.setInt("timeout", timeout)
	p.setJSON("allowed_updates", allowed)

	return performRequest[[]tg.Update](b, "getUpdates", p)
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
	if p.Scope != nil {
		v.setJSON("scope", p.Scope)
	}
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

func (b *Bot) setDefaultAdminRights(rights *tg.ChatAdministratorRights, forChannels bool) error {
	p := params{}
	p.setJSON("rights", rights)
	p.setBool("for_channels", forChannels)
	return performRequestEmpty(b, "setMyDefaultAdministratorRights", p)
}

func (b *Bot) getDefaultAdminRights(forChannels bool) (*tg.ChatAdministratorRights, error) {
	p := params{}
	p.setBool("for_channels", forChannels)
	return performRequest[*tg.ChatAdministratorRights](b, "getMyDefaultAdministratorRights", p)
}
