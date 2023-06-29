package tgot

import (
	"io"
	"runtime"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

type baseContext interface {
	Ctx() Context
}

// BaseContext can infer any context type.
type BaseContext[c baseContext] interface {
	baseContext
	Child(string) c
}

// MakeContext creates new context.
//
// The context must be used in a separate goroutine because in case of fatal
// errors the bot will be closed and the goroutine will be terminated.
func (b *Bot) MakeContext(name string) Context {
	return Context{bot: b, name: name}
}

// Context type.
type Context struct {
	bot  *Bot
	name string
}

// Ctx returns Context.
func (c Context) Ctx() Context { return c }

// Child creates sub context.
func (c Context) Child(name string) Context {
	c.name += "::" + name
	return c
}

// GetMe returns basic information about the bot.
func (c Context) GetMe() (*tg.User, error) {
	return method[*tg.User](c, "getMe")
}

// GetUserPhotos returns a list of profile pictures for a user.
func (c Context) GetUserPhotos(userID int64) (*tg.UserProfilePhotos, error) {
	d := api.NewData().SetInt64("user_id", userID)
	return method[*tg.UserProfilePhotos](c, "getUserProfilePhotos", d)
}

// GetFile returns basic information about a file
// and prepares it for downloading.
func (c Context) GetFile(fileID string) (*tg.File, error) {
	d := api.NewData().Set("file_id", fileID)
	return method[*tg.File](c, "getFile", d)
}

// DownloadReader downloads file as io.ReadCloser from Telegram servers.
func (c Context) DownloadReader(f *tg.File) (io.ReadCloser, error) {
	return c.bot.api.DownloadFile(f.FilePath)
}

// Download downloads file from Telegram servers.
func (c Context) Download(f *tg.File) ([]byte, error) {
	rc, err := c.DownloadReader(f)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// DownloadReaderFile downloads file as io.ReadCloser from Telegram servers
// by file id.
func (c Context) DownloadReaderFile(fid string) (io.ReadCloser, error) {
	f, err := c.GetFile(fid)
	if err != nil {
		return nil, err
	}
	return c.DownloadReader(f)
}

// DownloadFile downloads file from Telegram servers by file id.
func (c Context) DownloadFile(fid string) ([]byte, error) {
	f, err := c.GetFile(fid)
	if err != nil {
		return nil, err
	}
	return c.Download(f)
}

func (c Context) method(meth string, d ...*api.Data) error {
	_, err := method[api.Empty](c, meth, d...)
	return err
}

func method[T any](c Context, method string, d ...*api.Data) (T, error) {
	var data *api.Data
	if len(d) > 0 {
		data = d[0]
	}
	result, err := api.Request[T](c.bot.api, method, data)
	if err != nil {
		return result, c.bot.onError(c, result, err)
	}

	return result, err
}

// OnErrorDefault represents default error handler.
//
// It:
//
// - panics if the JSON response cannot be parsed;
//
// - panics if the method is not found;
//
// - cancels the current call (using runtime.Goexit) and
// changes the state of the bot to an error if 401 (Not Authorized) or
// or 500 (Internal Server Error) status is returned;
//
// - returns the err in any other case.
func OnErrorDefault(c Context, result any, err error) error {
	if e, ok := err.(*api.JSONError); ok {
		panic("incorrect Telegram JSON response (" + e.Method + "): " + e.Error())
	}

	tgErr, ok := err.(*api.Error)
	if !ok {
		return err
	}

	switch tgErr.Err.Code {
	case 404: // Not Found
		panic("Telegram method is not found but is present in the current implementation: " + tgErr.Method)
	case 401, 500: // Not Authorized, Internal Server Error
		c.bot.SetError(err)
		runtime.Goexit()
	}

	return err
}
