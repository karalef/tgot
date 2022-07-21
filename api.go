package bot

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
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

func (p params) set(key string, value interface{}) params {
	if value == nil {
		return p
	}
	vals := url.Values(p)
	switch v := value.(type) {
	case string:
		if v != "" {
			vals.Set(key, v)
		}
	case int:
		if v != 0 {
			vals.Set(key, strconv.Itoa(v))
		}
	case int64:
		if v != 0 {
			vals.Set(key, strconv.FormatInt(v, 10))
		}
	case bool:
		if v {
			vals.Set(key, strconv.FormatBool(v))
		}
	case float64:
		if v != 0 {
			vals.Set(key, strconv.FormatFloat(v, 'f', 6, 64))
		}
	default:
		b, _ := json.Marshal(value)
		vals.Set(key, string(b))
	}
	return p
}

func (b *Bot) requestURL(method string) string {
	return b.apiURL + b.token + "/" + method
}

func performRequest[T any](b *Bot, method string, p params) (T, error) {
	return performRequestContext[T](context.Background(), b, method, p)
}

func performRequestEmpty(b *Bot, method string, p params) error {
	_, err := performRequest[internal.Empty](b, method, p)
	return err
}

func decodeResponse[T any](method string, resp *http.Response, err error) (T, error) {
	var nilResult T
	switch err {
	case nil:
	case context.Canceled, context.DeadlineExceeded:
		return nilResult, err
	default:
		return nilResult, &Error{Err: err}
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
		return nilResult, r.APIError
	}

	return r.Result, nil
}

func performRequestContext[T any](ctx context.Context, b *Bot, method string, p params) (T, error) {
	data := url.Values(p)
	resp, err := internal.PostFormContext(ctx, b.client, b.requestURL(method), data)
	return decodeResponse[T](method, resp, err)
}

// File contains field and file data.
type File struct {
	Field string
	*tg.InputFile
}

func uploadFiles[T any](b *Bot, method string, p params, files []File) (T, error) {
	r, w := io.Pipe()
	mp := multipart.NewWriter(w)
	go func() {
		defer mp.Close()
		defer w.Close()

		err := p.forEach(func(field, value string) error {
			return mp.WriteField(field, value)
		})
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

	resp, err := b.client.Post(b.requestURL(method), mp.FormDataContentType(), r)
	return decodeResponse[T](method, resp, err)
}

func (b *Bot) downloadFile(path string) ([]byte, error) {
	resp, err := b.client.Get(b.fileURL + b.token + "/" + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (b *Bot) getMe() (*tg.User, error) {
	return performRequest[*tg.User](b, "getMe", nil)
}

func (b *Bot) getUpdates(ctx context.Context, offset, timeout int, allowed ...string) ([]tg.Update, error) {
	p := params{}
	p.set("offset", offset)
	//p.set("limit", limit)
	p.set("timeout", timeout)
	p.set("allowed_updates", allowed)

	return performRequestContext[[]tg.Update](ctx, b, "getUpdates", p)
}

type commandParams struct {
	Commands []tg.Command
	Scope    *tg.CommandScope
	Lang     string
}

func (p *commandParams) params() params {
	if p == nil {
		return nil
	}
	v := params{}
	v.set("language_code", p.Lang)
	if p.Scope != nil {
		v.set("scope", p.Scope)
	}
	v.set("commands", p.Commands)
	return v
}

func (b *Bot) getCommands(s *tg.CommandScope, lang string) ([]tg.Command, error) {
	p := commandParams{
		Scope: s,
		Lang:  lang,
	}
	return performRequest[[]tg.Command](b, "getMyCommands", p.params())
}

func (b *Bot) setCommands(p *commandParams) error {
	return performRequestEmpty(b, "setMyCommands", p.params())
}

func (b *Bot) deleteCommands(s *tg.CommandScope, lang string) error {
	p := commandParams{
		Scope: s,
		Lang:  lang,
	}
	return performRequestEmpty(b, "deleteMyCommands", p.params())
}

func (b *Bot) setDefaultAdminRights(rights *tg.ChatAdministratorRights, forChannels bool) error {
	p := params{}
	p.set("rights", rights)
	p.set("for_channels", forChannels)
	return performRequestEmpty(b, "setMyDefaultAdministratorRights", p)
}

func (b *Bot) getDefaultAdminRights(forChannels bool) (*tg.ChatAdministratorRights, error) {
	p := params{}
	p.set("for_channels", forChannels)
	return performRequest[*tg.ChatAdministratorRights](b, "getMyDefaultAdministratorRights", p)
}
